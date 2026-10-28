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
| Account entry | SHA-1, semantic digest serialization, and the dominant copy, projection, schedule, and encoding loops are recovered. Stage orchestration still uses the legacy CFG. |
| Bootstrap | Montgomery arithmetic and the dominant table/byte loops are recovered. Stage orchestration still uses the legacy CFG. |
| Selection | The 192-byte body and 224-byte envelope are parsed structurally. Candidate state/frame transforms, SIDB-key production, the 59-block SIDB AES consumer, and authenticated candidate promotion are typed Go. Native selection orchestration is no longer executed by model recovery. |
| Enrollment frame | The independent-key AES-shaped core and its separate nine-round suffix are normalized to a typed ten-layer AES-S-box SPN with affine 128-bit stages. Production no longer executes its generalized table model. |
| SIDB | Container projection, candidate-frame-to-key production, and the 59-block AES consumer are typed. The frame output encoding cancels the producer input encoding; ten AES-S-box stages plus affine 128-bit maps now derive each candidate's AES key directly. The legacy producer and per-enrollment DFA are no longer semantic requirements. |
| SINF | Portable AES/model decoding exists. Decoded SIDB discovery, payload-to-`valu` arithmetic, the fixed secondary operand, AES class join, exhaustive record/window evaluation, and unique personalized-SINF authentication are typed Go. Model recovery stops after selection; the native selector, SINF-entry CFG, and stack-offset material adapter are no longer executed. |
| SUPF | Variant-3 and variant-5 envelopes, architecture/page-key geometry, leaf certificates, and signed prefixes are parsed and authenticated in typed Go. Production pins the Apple FairPlay CA public key and no longer trusts a self-signed license certificate. |
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

The preceding `valu` producer is typed too. It consumes one 16-byte window
from the record's 20-byte `Payload || Tail` material. Each byte is XORed with
the corresponding position in one fixed 16-byte mask, wrapping when the
window crosses the payload boundary. Each resulting word then uses one
modular affine operation:

```text
semantic_word = 0xa46fd0ff * masked_window_word + 0xfe063531
                modulo 2^32
```

The former four apparent enrollment-specific word masks were a false boundary:
they are only rotations of the same fixed mask caused by different window
starts. The old layout uses offset 0; the independent current layout uses
offset 3 and naturally crosses into `Tail`. Both reproduce the exact protected
word entering the native arithmetic. The intermediate
XOR/add/multiply/mask MBA was removed. Its key identities were proved with
32-bit Z3 bit-vectors, and the direct operation matches the unreduced protected
arithmetic over 4,096 deterministic controls.

The secondary class operand is typed too. It contains one eight-byte value
twice, and its protected byte representation is only XOR `0x1b`. Unequal
halves are rejected. Adversarial validation corrected an earlier boundary
name: this value is the same build-global constant
`4ac6fcf2d3925e05` in both independent accepted enrollments. It decrypts each
account's different `valu` block to that account's already-proven model key
exactly. It is not produced by the account digest.

The digest-derived 3x3 workspace branch is real but separate from the class
secondary. Its arithmetic remains normalized in typed Go: three residues
modulo `32767`, anti-diagonal polynomial packing, one affine circulant map over
`GF(2^8)`, a 45th-power distance factor, the drift
`0x39a5 * 5^cell + 0x0c7e mod 32767`, and matrix cells
`(seed + index)^1000 mod 32767`. The seed comes from the final big-endian
16-bit account-digest suffix. The complete protected suffix domain matches a
direct reduction plus one, with eight bounded-carry exceptions. This branch
must not be used as the class key without a new dataflow proof.

The account digest itself is retained at a typed boundary. The native SINF
path previously serialized five SHA-1 state words through five protected loop
visits. One typed stage now writes the complete 20-byte digest under XOR
`0x1b` and retains its semantic value. A synthetic differential control
matches the complete protected loop.

The decoded current SIDB registry is now parsed structurally as fixed 40-byte
records followed by a 24-byte footer. Each record contains 20 bytes of
opaque native selection/validation metadata, a 16-byte `valu` payload, and a
four-byte tail whose semantic role remains open. Both independent decoded registries are
944 bytes: 23 records plus the footer, and the naturally selected record is
aligned within that array. Production obtains each registry directly from the
typed SIDB candidate decryptor, deduplicates its contents, and does not retain
a native allocation address or record offset.

Given the build-global class value, Go evaluates the typed `valu` and
AES class stages for every structurally parsed SIDB record and all five valid
record windows. Candidates retain their record, metadata, and window identity.
The installed personalized SINF then authenticates exactly one precursor by
decoding its private atom and checking the complete clear atom sequence. This
has replaced the native selector and the stack-offset SINF-material adapter in
model recovery. On the independent 26B118 control, 20 bounded registry-shaped
objects collapsed to one record/window candidate, produced a byte-identical
enrollment model with SHA-256
`78048f583dec3637404c24a9ed7bc6ae3313860c75207e81b8fabe3a45dc24b8`,
and decrypted 39 images and 10,080 pages. The decrypted IPA SHA-256 remained
`73fb374a683c2b10558ef8ab82b5833bdfab719a49a147597eb6dedf40fa325a`;
no extracted fallback block ran.

Model recovery now stops immediately after the outer stage. A current
224-byte source is parsed as a 32-byte envelope prefix plus a 192-byte body,
while the complete 224-byte layout remains a bounded alternative for two
additional slots. Typed Go evaluates every state/frame/SIDB-key candidate,
decrypts each current 993-byte SIDB container, and lets the personalized SINF
promote exactly one registry/record/window tuple. On the same private control,
26 typed registry candidates authenticated to the same single precursor,
produced the byte-identical model above, and completed the same IPA
decryption. Neither the generated selection CFG nor the generated SINF-entry
CFG is executed by model recovery.

SUPF authentication is typed and structural too. Variant 3 authenticates its
signed architecture/page-key prefix with RSA/SHA-1, while variant 5 uses
RSA/SHA-256. The embedded 1024-bit leaf certificate is checked against the
pinned Apple FairPlay Certification Authority public key before its key is
allowed to verify the prefix. The authority key was independently reconstructed
from three valid Apple certificate/signature pairs, reduced to one 1024-bit
modulus, and then verified against all ten FairPlay leaf certificates embedded
in the reference binary plus four v3/v5 licenses from two enrollment layouts.
The independent 26B118 private control still produced the byte-identical model
and decrypted IPA after this check was enabled.

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

1. Express account entry as typed stage functions rather than register-state
   adapters, preserving structural certificate, SIDB, and hardware inputs.
2. Express bootstrap as typed stage functions rather than register-state
   adapters.
3. Discover every remaining required input from SIDB/SINF/object structure,
   removing the legacy outer-transform and source/frame memory offsets.
4. Delete generated CFG and compact-interpreter production paths after typed
   stage parity is proven.
5. Canonicalize and compare the normalized algorithms across at least five
   independent builds and multiple enrollment layouts.

Update this file whenever a stage becomes typed, an offset is removed, a
cross-build invariant is proved, or a completion requirement changes state.
