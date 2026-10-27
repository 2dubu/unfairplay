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
| SIDB | Container projection and the 59-block AES consumer are typed. The legacy path was narrowed to an affinely encoded AES round-10 key; affine decoding plus ordinary AES-128 schedule reversal is now typed Go. The frame-to-encoded-round-key producer is still inside the legacy path. |
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

1. Recover the remaining typed candidate-frame-to-encoded-round-key producer.
   The encoded round key to SIDB decoder key stage is complete and must not be
   replaced by another DFA implementation.
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
