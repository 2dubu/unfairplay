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
| Enrollment frame | The independent-key AES-shaped core and its separate nine-round suffix are normalized to a typed ten-layer AES-S-box SPN with affine 128-bit stages. Production no longer executes its generalized table model. |
| SIDB | Container projection, candidate-frame-to-key production, and the 59-block AES consumer are typed. The frame output encoding cancels the producer input encoding; ten AES-S-box stages plus affine 128-bit maps now derive each candidate's AES key directly. The legacy producer and per-enrollment DFA are no longer semantic requirements. |
| SINF | Portable AES/model decoding exists. The 16-byte registry payload-to-`valu` arithmetic and downstream two-input class consumer are typed Go. The enclosing record still has to supply four word masks and the secondary class key without native orchestration. |
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

All nine generalized SIDB-key producer rounds are now separated from their
protected register representation. At each boundary, four complete 32-bit carriers were
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

These native-boundary rounds are not the enrollment-frame suffix described
below. They sit later, between the selected enrollment frame and the SIDB AES
key. Their protected boundary encodings are now removed too. The enrollment
frame's sixteen arbitrary output encodings compose with the producer prefix's
sixteen arbitrary input substitutions to sixteen affine byte maps. The nine
middle rounds and the final key layer are all AES-S-box layers separated by
affine 128-bit stages. The resulting producer profile is 43,392 bytes before
gzip and 1,929 bytes compressed.

The typed producer matched all one-byte and deterministic multi-byte controls
at every recovered boundary, then reproduced 17 independently DFA-recovered
AES keys exactly. A selected frame captured on the independent 26B118 machine
also produced the same key SHA-256 as a fresh full DFA recovery. Production
candidate enumeration now returns the state, portable frame, and derived SIDB
key. No handler address, register carrier, stack slot, generated CFG, or DFA
participates in that derivation.

The downstream class join is no longer an opaque white box. Its exact natural
path was isolated to 60 protected blocks. Coherent faults at its round-8,
round-7, and round-6 state boundaries completed the independent AES recovery
in about four seconds. The external input encoding is identity and the output
encoding is exactly XOR `0x1b` at every byte. After removing that output mask,
the result is exactly `AES-128-ECB-decrypt(class_key, valu_ciphertext)`.
Production now has this operation as a typed Go function; the protected core
must not be shipped or recovered again.

The preceding `valu` producer is typed too. Each four-byte payload word is
read big-endian and combined with one enrollment-supplied word mask. The
entire remaining protected arithmetic, byte codec, and following
representation mask reduce to one modular affine operation:

```text
semantic_word = 0xa46fd0ff * (payload_word XOR word_mask) + 0xfe063531
                modulo 2^32
```

The intermediate XOR/add/multiply/mask MBA was removed. Its key identities
were proved with 32-bit Z3 bit-vectors, and the direct operation matches the
unreduced protected arithmetic over 4,096 deterministic controls.

Together these stages reduce the remaining SINF seam to the enclosing
registry record's parameter producer. The secondary class operand is now typed
too: it contains one eight-byte enrollment value twice, and its protected byte
representation is only XOR `0x1b`. Unequal halves are rejected. This exact
representation matched both independent accepted enrollments. The remaining
producer must therefore yield four 32-bit word masks and one eight-byte class
value. Those values vary across the enrollments, while the payload arithmetic,
AES join, operand representation, and final material mask remain identical.
The generated CFG currently constructs those parameters; it is the next
normalization target.

Separately, the enrollment-frame suffix table algebra is now closed. All 144 input substitutions in its
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

1. Parse decoded SIDB records structurally and recover the metadata producer
   for the four `valu` word masks and secondary class key. Feed those values
   into the typed payload transform and AES class join; do not route either
   stage back through the generated CFG.
2. Try each typed SIDB candidate and select the unique candidate accepted by
   the authenticated SIDB/SINF structure.
3. Connect the resulting precursor to personalized SINF validation and remove
   native selection when exactly one candidate reaches a valid SINF.
4. Express account entry and bootstrap as typed stage functions rather than
   register-state adapters.
5. Discover required inputs from SIDB/SINF/object structure, removing the
   remaining memory offsets.
6. Delete generated CFG and compact-interpreter production paths after typed
   stage parity is proven.
7. Canonicalize and compare the normalized algorithms across at least five
   independent builds and multiple enrollment layouts.

Update this file whenever a stage becomes typed, an offset is removed, a
cross-build invariant is proved, or a completion requirement changes state.
