# Build-independent recovery plan

## Goal

FairPlay enrollment must use one shared Go implementation. A macOS build must
not need translated assembly, a flattened CFG, handler addresses, register
scripts, or hard-coded stack offsets.

Different `fairplayd` binaries are treated as differently obfuscated versions
of the same algorithm unless dataflow evidence proves a real semantic change.

The final flow should be:

```text
SIDB + SINF + hardware identity
        -> typed Go enrollment algorithms
        -> portable enrollment model
        -> IPA decryption
```

## Design rules

- Algorithms accept typed values, slices, tables, and state structures.
- Algorithms do not accept instruction addresses or emulated ARM registers.
- Stage inputs are found by validated file/object structure, not build offsets.
- Build hashes are evidence labels, not algorithm selectors.
- A genuine changing table may be data. Changing CFGs, MBA constants,
  registers, selectors, and handler addresses are obfuscation noise.
- Unsupported input geometry fails closed.
- Native binaries, IDA, Binary Ninja, Python, and compact interpreters are
  research oracles only. They are not runtime dependencies.
- A new build must not require another generated CFG.

## Typed stage contracts

The intended stage boundary is:

```text
Account input
  SIDB records, hardware identity, installed SINF metadata

Bootstrap
  typed big-integer, SHA-1, table, and byte transformations

Selection
  192/224-byte enrollment source -> candidate 16-byte states and frames

SIDB validation
  candidate frame -> enrollment-bound SIDB decoder key
  SIDB container + decoder key -> decoded registry

SINF validation
  decoded SIDB registry + personalized SINF -> accepted precursor

Model output
  precursor -> current SINF key/tables -> portable enrollment model
```

When candidate validation can replace native selection, it should. Trying all
typed candidates against the authenticated SINF is preferable to reproducing a
flattened native selector.

## Current status

| Area | Status |
| --- | --- |
| Outer transform | Typed Go algorithm is recovered. Its legacy memory adapter still has offsets and must be removed. |
| Digest, pair, tail, writer | Typed Go algorithms are recovered and differential-tested. |
| Account entry | SHA-1 and the dominant copy, projection, schedule, and encoding loops are recovered. Stage orchestration still uses the legacy CFG. |
| Bootstrap | Montgomery arithmetic and the dominant table/byte loops are recovered. Stage orchestration still uses the legacy CFG. |
| Selection | Candidate state/frame transforms and several dominant loops are recovered. Native orchestration is still transitional. |
| SIDB | Container projection and the 59-block AES consumer are typed. The legacy path was narrowed to an affinely encoded AES round-10 key; affine decoding plus ordinary AES-128 schedule reversal is now typed Go. The nine generalized frame-producer rounds and their final shared-output layer are now normalized to a typed AES-S-box SPN with affine 128-bit stages. |
| SINF | Portable AES/model decoding exists. Directly treating a candidate state or frame as the final precursor was tested and rejected because SIDB registry decoding comes first. |
| Build portability | The portable algorithm layer is being separated from the legacy 26A428 oracle. Zero-offset object discovery is still incomplete. |

The current `executionProfile`, generated block files, emulated registers, and
semantic memory-layout adapters are transitional oracle infrastructure. They
are not the final architecture.

The prior per-enrollment SIDB DFA is no longer a semantic requirement. At the
pre-consumer boundary, sixteen encoded bytes decode to the complete AES-128
round-10 key. Reversing the standard key schedule reproduces the SIDB decoder
key. This was derived from 64 chosen-frame controls and then matched 17
independent recovered keys exactly. The affine byte encodings are still
26A428 oracle evidence until they are normalized across other builds.

The protected producer's latter half is also typed: one-byte differential
masks match the ordinary AES-128 schedule for rounds 6, 7, 8, and 9 exactly.
The Go schedule reversal now accepts any round from 0 through 10. No separate
round-specific white box or DFA is required after a semantic round key has
been recovered. The remaining unknown is the earlier frame-to-round-key
producer and its protected byte representation.

All nine generalized producer rounds are now separated from their protected
register representation. At each boundary, four complete 32-bit carriers were
identified while transient XOR shares were excluded. Every round is exactly a
baseline plus sixteen independent 8-to-32-bit XOR contribution tables, with
four input bytes feeding each four-byte output group. For every round, all
4,096 one-byte controls and 256 deterministic multi-byte controls matched the
typed `generalizedByteRound` shape. The independently recovered round digests
are:

```text
25b4b4e6 0ba5905a e896a7fe d2a1bf8b cdea5bf1
aa1f69d0 36867136 e797caf6 64556dd1
```

These native-boundary coordinates deliberately do not hash like the existing
semantic profile's internal rounds: the older trace solver chose different
byte permutations and protected boundary encodings. Their complete transforms
still agree at the public frame output.

The table algebra is now closed. All 144 input substitutions in the
nine generalized rounds are exactly affine-equivalent to the ordinary AES
S-box, exhaustively over the 256-byte domain. For each substitution, its four
output components are related by invertible affine byte maps. The final layer
contains four direct byte permutations and twelve shared outputs; all 36
nonlinear contributions to those shared outputs are also exactly
affine-equivalent to the AES S-box. The twelve shared outputs form four
four-byte groups. Each group reduces to an invertible 32x32 affine map followed
by four byte permutations. Composing those maps gives the complete suffix as
ten parallel AES-S-box layers separated by affine 128-bit stages, followed by
one affine output stage and sixteen byte encodings.

Production now executes that normalized SPN directly. Its checked profile is
51,592 bytes before gzip and 17,024 bytes compressed, down from the 167,844-byte
generalized table model. The old generalized model is test-only and serves as
an independent oracle. Across 4,096 deterministic controls, the normalized SPN
matches both its 16-byte output and its dispatcher exit word exactly. A direct
DarkPhoenix run was also a useful negative: it reaches step 2 and fails the
Tolhuizen condition because the unreduced final layer is groupwise rather than
bytewise. The group reduction explains that failure; another DFA run is not
needed.

## Cross-build evidence

Local binaries already show that CFG identity is not a useful definition of
build support:

- 26A428 and 26B118 have very different normalized anchors and flattened CFGs.
- The recovered SHA-1, Montgomery, byte-mask, word-mask, selection-state, and
  SINF algorithms remain valid semantic operations.
- Therefore cross-build work is used to remove obfuscation-dependent details,
  not to add one implementation per binary.

At least five independently hashed release, RC, and update builds must receive
stage-level differential validation before this work is complete.

## Proof requirements

Every promoted transformation needs suitable evidence:

- exhaustive comparison for byte-sized domains;
- randomized differential tests against the exact Go oracle;
- Z3 or equivalent bit-vector proof for MBA identities where practical;
- comparison across independent binary builds;
- exact private QA hashes, image count, and decrypted page count;
- `go test ./...`, `go vet ./...`, and `go test -race ./...`.

## Next work

1. Compare and canonicalize the normalized SPN affine stages across independent
   macOS builds. The 26A428 constants are semantic evidence, not a valid build
   selector. Do not add another generated CFG or generalized table profile for
   a new build.
2. Parse the decoded SIDB registry structurally and recover its typed
   registry-to-current-SINF-precursor transform.
3. Connect the precursor to personalized SINF validation and remove native
   selection when exactly one candidate reaches a valid SINF.
4. Express account entry and bootstrap as typed stage functions rather than
   register-state adapters.
5. Discover required inputs from SIDB/SINF/object structure, removing the
   remaining memory offsets.
6. Delete generated CFG and compact-interpreter production paths after typed
   stage parity is proven.
7. Validate the same algorithms across at least five independent builds and
   multiple enrollment layouts.

Update this file whenever a stage becomes typed, an offset is removed, a
cross-build invariant is proved, or a completion requirement changes state.
