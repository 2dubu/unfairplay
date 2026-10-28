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
| Outer transform | Typed Go constructs its complete input from validated 630-, 993-, or 1003-byte SIDB geometry and the six-byte hardware identity. The legacy stack/register adapter is no longer used by production recovery. |
| Digest, pair, tail, writer | Typed Go algorithms are recovered and differential-tested. |
| Account entry | Replaced by structural X.509 parsing and pinned-CA SUPF authentication. The legacy account-entry CFG has been removed. |
| Bootstrap | Replaced by one typed semantic outer profile plus SIDB/hardware runtime fields. The legacy bootstrap CFG has been removed. |
| Selection | The 192-byte body and 224-byte envelope are parsed structurally. Candidate state/key transforms, SIDB-key production, the 36- or 59-block SIDB AES consumer, and authenticated candidate promotion are typed Go. Native selection orchestration is no longer executed by model recovery. |
| Enrollment frame | The independent-key AES-shaped core and its suffix are normalized to direct keyed AES S-box transitions and sparse `GF(2^8)` mixing layers. The larger generalized and affine-SPN models are test-only. |
| SIDB | Container projection, candidate-state-to-key production, and the 36- or 59-block AES consumer are typed. The unused frame output encoding is gone. The producer is now ten direct AES inverse-S-box layers, nine fixed `InvShiftRows + InvMixColumns` operations, four final XOR graphs, and compact endpoint maps. Its protected stage encodings and per-enrollment DFA are no longer semantic requirements. |
| SINF | Portable AES/model decoding exists. Decoded SIDB discovery, current and older payload-to-`valu` arithmetic, all semantic class-value families, AES class join, exhaustive record/window or mask-phase evaluation, and unique personalized-SINF authentication are typed Go. The native selector, SINF-entry CFG, and stack-offset material adapter have been removed. |
| SUPF | Variant-3 and variant-5 envelopes, architecture/page-key geometry, leaf certificates, and signed prefixes are parsed and authenticated in typed Go. Production pins the Apple FairPlay CA public key and no longer trusts a self-signed license certificate. |
| Build portability | The 630-, 993-, and 1003-byte SIDB layouts recover through the zero-offset Go path without a native executable, reference binary, generated CFG, compact model, per-enrollment DFA, or iPhone-derived secret. The five-build account-entry static matrix is complete; structurally locating the downstream frame and SIDB-key stages in independent builds remains open. |

The legacy `executionProfile`, generated account/bootstrap/selection blocks,
emulated registers, compact interpreter, exact-build fallbacks, and native
memory-layout adapters have been removed. The production recovery package now
contains typed algorithms plus only the compact frame-AES and SIDB-key data
profiles.

Those production artifacts and symbols are build-neutral (`current-frame-aes`
and `current-sidb-producer-typed`). The larger 26A428-labeled generalized and SPN
models are test-only and remain solely as independent normalization oracles.

New captures contain only structural enrollment inputs. Collection no longer
copies fairplayd, profiles its CFG, or records a daemon hash as a recovery
requirement. Recovery no longer accepts a reference-binary override. Legacy
capture manifests that contain the old diagnostic binary remain readable.

New enrollment models also no longer identify themselves with the 26A428
fairplayd hash. Their 32-byte identity field is now the SHA-256 of the semantic
domain string `unfairplay/typed-current-enrollment/v1`:
`a3aa881749848636d88a8b0ec02adb35b77101d0a001dea6a0671696d8595402`.
Model selection does not compare this value with a capture or host daemon;
the selector and complete personalized SINF authentication are authoritative.
Existing version-1, version-2, and version-3 models remain readable. A retained
pre-change version-2 model still verified the same 10 images and 11,004 pages,
with decrypted-IPA SHA-256
`0eafcbac1502c38aeaa069a627e2d6592f0f2a19b42222361c392a85d77f79ea`.

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
frame's final affine map composes into the producer's first layer. All nine
middle layers factor exactly as AES `InvMixColumns`. A global gauge solve reduces
every intervening nonlinear transition to the AES inverse S-box with a direct
round-key XOR; every one of the 36 factor domains collapses from 1,020 candidates
to four equivalent gauges. Output masks compose into the following key or final
constant. A second byte-order solve maps all nine round matrices to one fixed
`InvShiftRows + InvMixColumns` operation; only 16 equivalent global orderings
remain. The eight recovered keys are independent round constants rather than
an AES-128 key schedule. The former final 40-map affine layer factors exactly
into four byte-wise XOR graphs. One unique gauge per graph turns the preceding
sixteen protected substitutions into direct keyed AES inverse S-boxes; their
output masks fold into the final constants. The initial protected endpoint has
one unique normalization per byte: all
sixteen substitutions are direct AES inverse S-boxes with a zero key. Their
output masks fold into the first middle-round key. The complete producer is
therefore ten inverse-S-box layers from end to end. The checked production
profile is 528 bytes before gzip and 390 bytes compressed. It contains the
round constants, four XOR adjacency masks, and compact endpoint maps. The
former 43,392-byte affine-SPN profile is now a test-only 26A428 normalization
oracle.

The typed producer matched all one-byte and deterministic multi-byte controls
at every recovered boundary, then reproduced 17 independently DFA-recovered
AES keys exactly. The inverse-AES form additionally matched the former
affine-SPN oracle over 4,096 deterministic controls. A selected frame captured
on the independent 26B118 machine also produced the same key SHA-256 as a fresh
full DFA recovery. Production candidate enumeration now returns the semantic
state and derived SIDB key; the unused protected frame encoding has been
removed. No handler address, register carrier, stack slot, generated CFG, or
DFA participates in that derivation.

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
halves are rejected. Current compact enrollments use one fixed semantic value.
The older extended enrollment uses the output of the digest-derived 3x3
workspace branch. That dataflow is now proved by an independent class-core
DFA: the workspace output, repeated twice as the AES key, decrypts the typed
record `valu` to the already-proven model key exactly.

The workspace arithmetic remains normalized in typed Go: three residues
modulo `32767`, anti-diagonal polynomial packing, one affine circulant map over
`GF(2^8)`, a 45th-power distance factor, the drift
`0x39a5 * 5^cell + 0x0c7e mod 32767`, and matrix cells
`(seed + index)^1000 mod 32767`. The seed comes from the final big-endian
16-bit account-digest suffix. The complete protected suffix domain matches a
direct reduction plus one, with eight bounded-carry exceptions. This branch
now supplies the older class-value family.

Production does not reconstruct the native account-digest path. Its only
downstream effect is the workspace seed, whose complete 32,767-value semantic
domain is cheap to enumerate. Personalized-SINF authentication selects the
single coherent value, so the old SHA-1 writer and its CFG are unnecessary.

The older 630-byte enrollment family is closed separately without that seed
enumeration. Its semantic account digest is ordinary SHA-1 over a typed
45-byte message: a fixed 17-byte prefix, the six-byte built-in IOMACAddress
plus two zero bytes, one reserved zero byte, the one-byte SINF key selector,
the big-endian SINF user selector, two `0xff` bytes, and a fixed 12-byte
suffix. The final two big-endian digest pairs are reduced modulo 32767 and
each drives one direct 3x3 coordinate/velocity recurrence. The two resulting
polynomial words are combined by one circulant `GF(2^8)` map with column
`e1 2e 9b b9 e8 6d ff df`; the result is repeated to form the AES class key.
The native three-matrix limb encodings, word biases, byte tables, and final
XOR-`0x1b` transport have all been removed.

The older `valu` path is the same four-word modular affine transform used by
current records. It consumes the complete 16-byte payload and starts the fixed
record mask at phase 8. Production evaluates all sixteen bounded mask phases
and lets complete personalized-SINF authentication select the unique coherent
record and phase, so phase selection does not depend on a native offset or
dispatcher.

The independent older private control recovered from its 630-byte Mac SIDB,
IOMACAddress, installed license, and personalized IPA with no iPhone or iOS 16
material. The new typed model SHA-256 was
`6f9888452b6f9f0acf5466d32b71ee52f6b4ddca3d1e27027879b4f18336ef74`.
It decrypted one encrypted image and one page; decrypted-IPA SHA-256 was
`d71c1a150acc334b0dc0bf9532bc8501aa8c2020cb7d906b9b4cb5b212ce0f53`,
and logical-key SHA-256 matched the established independent result
`0cc0921e90789a30c7cce2ba80c75ced64dd334af31ff6744dc714923a14ac32`.

The decoded current SIDB registry is now parsed structurally as fixed 40-byte
records followed by a 24-byte footer. Each record contains 20 bytes of
opaque native selection/validation metadata, a 16-byte `valu` payload, and a
four-byte tail whose semantic role remains open. Both independent decoded registries are
944 bytes: 23 records plus the footer, and the naturally selected record is
aligned within that array. Production obtains each registry directly from the
typed SIDB candidate decryptor, deduplicates its contents, and does not retain
a native allocation address or record offset.

The current-SINF external codec is identity after normalization. Its protected
form XORed each byte with a position mask and then used an inverse table that
XORed the same mask again. Production now copies the authenticated private
source directly and emits zero whitening plus canonical identity tables for
backward-compatible version-2 models. The former semantic-profile object and
its 4,112 bytes of encoding data have been removed.

Go evaluates the fixed current class value and every value in the complete
32,767-element workspace-seed domain against the typed `valu` and AES class
stages for every structurally parsed SIDB record and all five valid record
windows. A one-block SINF prefix probe makes the exhaustive semantic search
bounded and fast; the complete SINF is still authenticated before promotion.
Candidates retain their record, metadata, and window identity.
The installed personalized SINF then authenticates exactly one precursor by
decoding its private atom and checking the complete clear atom sequence. This
has replaced the native selector and the stack-offset SINF-material adapter in
model recovery. On the independent 26B118 control, 20 bounded registry-shaped
objects collapsed to one record/window candidate, produced a byte-identical
enrollment model with SHA-256
`98c29933c57e3425f22005bfaf43fac684ef6b671629f76361cb1dd502147575`,
and decrypted 39 images and 10,080 pages. The decrypted IPA SHA-256 remained
`73fb374a683c2b10558ef8ab82b5833bdfab719a49a147597eb6dedf40fa325a`;
no extracted fallback block ran.

Model recovery now stops immediately after the outer stage. A current
224-byte source is parsed as a 32-byte envelope prefix plus a 192-byte body,
while the complete 224-byte layout remains a bounded alternative for two
additional slots. Typed Go evaluates every state/frame/SIDB-key candidate,
decrypts each supported SIDB container, and lets the personalized SINF
promote exactly one registry/record/window tuple. On the same private control,
26 typed registry candidates authenticated to the same single precursor,
produced the byte-identical model above, and completed the same IPA
decryption. The generated selection and SINF-entry CFGs have been removed.

The older container is typed as a 43-byte header, a 944-byte encrypted
registry, and a 16-byte trailer. Its 39-byte codec drives 427 outer iterations;
the registry uses the same typed SIDB input, output, and chaining masks as the
compact layout. On the independent 26A434 control, the typed path produced a
byte-identical enrollment model with SHA-256
`45c6529775ecb96e5bb3449324ab5b751e2fdecf1699724784007ec93fde372a`,
selected one authenticated record/window from 26 decoded registry candidates,
and decrypted 10 images and 11,004 pages. The current decryption output had
SHA-256 `0eafcbac1502c38aeaa069a627e2d6592f0f2a19b42222361c392a85d77f79ea`;
all 1,024 audited ARM64 words on the selected page were valid. No native or
extracted block ran.

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

The account-entry and bootstrap boundary has now collapsed too. A full live-in
audit found 3,671 bytes crossing from account entry to bootstrap, mostly native
X.509 objects, Montgomery scratch, descriptors, and pointers. Invalid signed
prefix or signature mutations both take the native rejection edge. Four valid
v5 licenses from e-Devlet, 9292, WhatsApp, and Airalo—including compact and
multi-page layouts—produce the exact same 167-word outer table, five state
words, optional word, and loop controls. A hardware mutation changes none of
those values; it changes only the explicit six-byte hardware input. The compact
29-byte and extended 39-byte codecs are exactly the bytes between the four-byte
version and their structurally located encrypted registry. Production
therefore verifies the SUPF in ordinary Go, constructs
the fixed semantic outer profile, copies the SIDB codec and hardware bytes, and
runs the typed outer algorithm directly. The old `fairplayd`, checkpoint,
account/bootstrap CFG, and stack-offset adapter dependencies have been removed.

The per-codec outer cycle is no longer a generated MBA circuit in production.
Its five 32-bit state updates are direct XORs of rotated codec and hardware
bytes. Its five wide continuation values reduce to the hardware index, two
selected hardware bytes, and two small affine index expressions. The direct
implementation matched the retained 764-line protected circuit exactly on
4,096 deterministic controls; that circuit now exists only as a test oracle.

On the independent 26B118 current-layout control, this zero-offset path still
produced byte-identical enrollment-model SHA-256
`98c29933c57e3425f22005bfaf43fac684ef6b671629f76361cb1dd502147575`,
decrypted 39 images and 10,080 pages, and retained decrypted-IPA SHA-256
`73fb374a683c2b10558ef8ab82b5833bdfab719a49a147597eb6dedf40fa325a`.
No native or extracted block ran.

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

The remaining affine gauges are solved globally rather than selected one layer
at a time. Each of the 36 four-byte components initially has 1,020 exact AES
factorizations. Constraints at every shared byte boundary reduce every domain
to four equivalent gauges. A coherent choice exposes eight direct keyed AES
S-box transitions separated by nine sparse `GF(2^8)` mixing layers. The four
apparent output masks are fixed coordinate encodings, not new nonlinear
operations.

Production now executes that compact AES-S-box/MDS form directly. Its checked
profile is 4,940 bytes before gzip and 2,065 bytes compressed, down from the
51,592-byte affine SPN and the original 167,844-byte generalized table model.
The frame's unused diagnostic output encoding and protected dispatcher exit
word are no longer production state. Its final affine output is composed into
the SIDB-key producer bridge. The generalized and affine SPN models are
test-only independent oracles. Across 4,096 deterministic controls, the
compact form reproduces the exact affine-SPN state. A direct DarkPhoenix run
was also a useful negative: it reaches step 2 and fails the Tolhuizen condition
because the unreduced final layer is groupwise rather than bytewise. The gauge
solve explains that failure; another DFA run is not needed.

## Cross-build evidence

Local binaries already show that CFG identity is not a useful definition of
build support:

- The stable account-entry control matrix now spans five independently hashed
  releases:

  | Build | arm64e SHA-256 | `__text` SHA-256 | Words identical to 26A428 |
  | --- | --- | --- | ---: |
  | 22A380 | `b09ad07348949aa0dbcfad8058cf1871a95e45caed575c58db9f9c78ba600c66` | `9d9637c9c5e0e91c886a494b3c7c08c9554156b145992b4c1210ae06f3eb4e85` | 9,710 / 9,779 |
  | 23A344 | `07e2b83fbeee23df0a2262188d55f9bda19f0c0dec0a78b5839b74a6b940f546` | `d53d844510c676bf002671d6543221428dc9789d4065357b0db3a32ed7a92dfa` | 9,684 / 9,779 |
  | 24A335 | `df4986011a9dff6f8078dc99d9d0cd27d8fd0af248e2034145b30ebd6041d18b` | `466b359399c23243e21387465e2f1ae96e8870014b88bfdf92345e2b7f89cff3` | 9,684 / 9,779 |
  | 25G83 | `93b8db9c7140cedbd04a3bbbce297b68ba9060e8a4f39ef52a1017964f3da4c5` | `425d914b6164320bb06fb9ae891d522dd0aeab719ca577825f8076490be30f7e` | 9,684 / 9,779 |
  | 26A428 | `1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784` | `eb9214dd24610e71505223d41c7627575952de101e44d474dd2fc06bea9045a4` | reference |

  Every complete account-entry function has 9,779 instructions and normalized
  SHA-256
  `b0dc6ba4f1e4f57a4e9385c65774dfa5bf27040e4ee9f60c3e3128b10f463da2`.
  For 23A344 through 25G83, all 95 non-identical words are relocated
  references: 47 `ADRP`, 29 `LDR`, 18 `ADD`, and one `BL`. For 22A380, the
  corresponding counts are 21, 29, 18, and one. No arithmetic or control
  instruction differs. This is full static differential evidence for the
  account-entry implementation across four major-version boundaries. It does
  not yet prove equal runtime outputs or validate the downstream frame and
  SIDB-key stages.
- The recovered SHA-1, Montgomery, byte-mask, word-mask, selection-state, and
  SINF algorithms remain valid semantic operations.
- Therefore cross-build work is used to remove obfuscation-dependent details,
  not to add one implementation per binary.

The former downstream BinDiff anchor map is rejected. Its one-block matches
were mnemonic-shape false positives, not typed dataflow matches. The clearest
control is the proposed 25G83 frame candidate at `0x100a6ef30`: BinDiff paired
it with 26A428 `0x1002fbb24` at similarity `0.428845` and confidence
`0.622459`, but the correspondence fails at the semantic boundary. Where the
26A428 path reads the sixteen input bytes, 25G83 performs 32-bit shifts and
bit slicing into a stack word schedule. The proposed 25G83 path continues for
4,136 instructions into a multiword compression-shaped loop and never exposes
the reference function's sixteen-byte input/output geometry. Matching mnemonic
counts therefore do not locate a protected stage. The associated 22A380,
23A344, and 24A335 candidates are demoted too until structural dataflow proves
them independently.

Cross-build discovery must start from semantic invariants: one structurally
identified sixteen-byte input, one sixteen-byte output, complete causal
diffusion, and normalization to the typed Go stage. Function similarity,
address proximity, and flattened-handler roots may rank candidates, but cannot
promote one.

That structural search has now recovered a narrower, evidence-ranked set:

- The proven 26A428 suffix at `0x1002fbb24..0x1002ff8d8` and the 25G83 path at
  `0x100086a70..0x10008a834` both contain exactly 3,953 instructions. The
  25G83 path has the same complete sixteen-byte input/output geometry, a
  mnemonic-histogram overlap of 3,905/3,953, and a mnemonic sequence ratio of
  `0.736403`. This is a credible local frame-family candidate, but it is not a
  proved natural-path match and its semantic output still needs normalization
  against the typed Go stage.
- The 22A380, 23A344, and 24A335 images each contain roughly eighty protected
  functions with the same ten-round shape: 160 indexed 32-bit table loads,
  sixteen-byte input and output geometry, and complete causal diffusion.
  Representative members at 22A380 `0x100c1c824`, 23A344 `0x10071bc74`, and
  24A335 `0x1001a0704` execute deterministically in the exact offline images.
  This proves a shared white-box family, not that any representative is the
  enrollment-frame instance. Apple includes many alternate or decoy instances,
  so diffusion and table counts cannot select the semantic path.
- Pointer-table ranks, selector ordinals, and effective table-address order are
  randomized between 25G83 and 26A428. All three can drive valid-looking
  sixteen-byte transforms, but none normalizes to the typed frame core. They
  are therefore rejected as cross-build identities. The exact natural caller
  state must select the instance and its tables.
- Local predecessor matching is insufficient too. The exact 26A428 frame path
  enters through `0x100242ea8`, loops through `0x1003d0dd8`, and reaches its
  descriptor caller at `0x100433434`. The independently top-ranked 25G83
  lookalikes are `0x100030084`, `0x1000d9778`, and `0x100034d60`, but they do
  not form one coherent protected chain. Under the only dispatch-table base and
  entry state that send the first block to the proposed loop, that loop's next
  encoded index points outside the same static fixup table. This rejects the
  composed match; it does not prove a new algorithm. Cross-build discovery must
  preserve one coherent dispatcher state from entry through frame output.

The formerly top-ranked 22A380 candidate at
`0x10061d8c8..0x100623954` is rejected. It has attractive instruction counts,
sixteen-byte reads, and a `0.927677` mnemonic-histogram overlap, but changing
each apparent input byte caused zero output changes. This is the concrete
control showing why geometry and mnemonic similarity are insufficient without
causal mutation.

These are research controls, not build-specific production profiles. Promotion
requires showing that an exact natural path decodes to the same semantic Go
operation. The addresses, table order, selector values, register allocation,
and dispatcher encodings must not enter production.

A separate 64-byte-state plus 16-byte-input helper family is explicitly not
the frame. It appears at 22A380 `0x100b2e8a0`, 24A335 `0x100efd024`, and
25G83 `0x100bf2928`, and produces byte-identical chosen-control outputs across
those builds. Its sixteen 32-bit block loads and compression-shaped dataflow
make it a useful cross-build control, but its 64+16-to-16 contract is the wrong
semantic boundary. Do not promote it as enrollment evidence.

At least five independently hashed stable release and update builds must
receive stage-level differential validation before this work is complete.
The cross-generation control set starts with macOS 13.0 (22A380), macOS 14.0
(23A344), and macOS 15.0 (24A335), then adds two stable macOS 26/27 builds.
These binaries are research controls: inclusion here is not a claim that their
enrollment layouts are already supported.

## Proof requirements

Every promoted transformation needs suitable evidence:

- exhaustive comparison for byte-sized domains;
- randomized differential tests against the exact Go oracle;
- Z3 or equivalent bit-vector proof for MBA identities where practical;
- comparison across independent binary builds;
- exact private QA hashes, image count, and decrypted page count;
- `go test ./...`, `go vet ./...`, and `go test -race ./...`.

## Next work

1. Structurally locate the true frame and SIDB-key boundaries in the independent
   builds. First recover one coherent 25G83 dispatcher path from entry through
   the descriptor caller and normalize its frame against the typed Go stage;
   do not compose independently ranked lookalike blocks. Use the resulting
   semantic fingerprint to select the enrollment member from the many 22A380,
   23A344, and 24A335 white-box-family controls. Require sixteen-byte
   input/output geometry, causal dataflow, coherent control state, and exact
   semantic agreement; compare outputs, not protected registers, instruction
   words, selector ranks, or table addresses. Classify any genuine constant
   change as input data rather than a build selector.
2. Run exact private QA across multiple enrollment layouts after each promoted
   cross-build normalization. The three older binaries are differential
   controls and do not require private account enrollment on those OS releases.

Update this file whenever a stage becomes typed, an offset is removed, a
cross-build invariant is proved, or a completion requirement changes state.
