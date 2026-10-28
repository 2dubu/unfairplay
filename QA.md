# Private QA

The control matrix is manual and opt-in. Private captures are never stored in
the repository. Temporary decrypted IPAs are deleted when the command exits.

Build UnfairPlay, create a mode-0600 manifest, then run:

```sh
go build -o /private/tmp/unfairplay ./cmd/unfairplay
go run ./tools/qa-controls \
  --unfairplay /private/tmp/unfairplay \
  --manifest /private/path/controls.json
```

Example manifest:

```json
{
  "schema": "unfairplay-private-qa-v1",
  "controls": [
    {
      "name": "modern-control",
      "capture": "/private/path/App.capture",
      "expected_model_sha256": "...",
      "expected_ipa_sha256": "...",
      "expected_reference_code_source": "typed-go",
      "expected_images": 10,
      "expected_pages": 11004
    },
    {
      "name": "saved-ipa-control",
      "ipa": "/private/path/App.ipa",
      "model": "/private/path/fairplay-enrollment.bin",
      "expected_model_sha256": "...",
      "expected_ipa_sha256": "...",
      "expected_reference_code_source": "existing-model",
      "expected_images": 39,
      "expected_pages": 10080
    }
  ]
}
```

Capture controls test complete recovery. IPA/model controls test reusable
decryption across other accounts and app layouts. The command fails on any
changed hash, image count, page count, or incomplete decryption. Use it before
creating a release tag.

## macOS IPA conversion control

The private conversion control proves that a version-2 enrollment model derives
the same account credential and output SINF as the independent macOSFairFree
implementation:

```sh
UNFAIRPLAY_MACOS_CREDENTIAL_MODEL=/private/path/fairplay-enrollment.bin \
UNFAIRPLAY_MACOS_CREDENTIAL_IPA=/private/path/App.ipa \
go test ./internal/fairplay \
  -run TestExternalMacOSCredentialFromEnrollmentModel -count=1 -v
```

Use an IPA containing exactly one current 448-byte private SINF for this
control. Do not commit the model, IPA, SINF, credential, or converted output.

## Compact frame AES profile

The production frame suffix uses direct keyed AES S-box transitions and nine
sparse `GF(2^8)` mixing layers. The larger affine SPN is retained only as an
independent test oracle. Verify all 4,096 deterministic controls with:

```sh
go test ./internal/recovery -run '^TestFrameAESProfile$' -count=1
```

Expected size: 2,065 bytes. Expected compressed SHA-256:
`9d42fb362c39a1f5e1c152fb82fa371869afa0701bec3065b89fba10915af202`.
The decoded checked profile is 4,940 bytes including its checksum, with
SHA-256 `33b7a8b3d13d7d5d09d665f52f5f887c426bd04e119b604d09136585b2ae0945`.
The test-only affine-SPN oracle remains 17,024 bytes with SHA-256
`f1cb13c0056a4125387f98a575c21cc6129438580e48c76373e333df614ea963`.

## Typed inverse-AES/XOR SIDB-key producer

Production composes away the protected producer's separate stage encodings.
It is ten direct AES inverse-S-box layers, nine fixed
`InvShiftRows + InvMixColumns` operations, four final XOR graphs, and compact
input/output maps. The first inverse-S-box key is zero; all of its former masks
fold into the next round key. Per-round byte orders and matrices are gone.
Verify all 4,096 deterministic controls against the larger test-only SPN with:

```sh
go test ./internal/recovery -run '^TestSIDBProducerTypedMatchesSPNOracle$' -count=1
```

Expected size: 390 bytes. Expected compressed SHA-256:
`a1ddefeea78738a85699db79f2a2c1687b87201c69f43db944e80cbd696ee5be`.
The decoded checked profile is 528 bytes including its checksum, with
SHA-256 `7c0b374de0559f5aa5e743e358b1247636100051e98172542256332ba421b4b9`.

The original private controls can regenerate the larger normalization oracle:

```sh
UNFAIRPLAY_SIDB_PRODUCER_PROFILE=/private/path/producer.bin.gz \
UNFAIRPLAY_SIDB_KEY_LAYER_PROFILE=/private/path/key-layer.bin \
UNFAIRPLAY_SIDB_KEY_REPORT_INPUT=/private/path/key-controls.json \
UNFAIRPLAY_SIDB_PRODUCER_SPN_OUTPUT="$PWD/internal/recovery/profiles/macos26a428-sidb-producer-spn-oracle.bin.gz" \
go test ./internal/recovery -run '^TestSIDBProducerSPNResearch$' -count=1
```

The test-only oracle is 1,618 bytes. Its compressed SHA-256 is
`7f31d8e3b36bcd2981ff2edd8d742d83616c925a9a5c744c868af8e576edc904`.
The 43,392-byte decoded profile has SHA-256
`f9362ea765534d3a3e578f5e40d3814d8cfa57dc67306aadf0e0e698b4344954`.
