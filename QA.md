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

## Normalized frame SPN profile

The production frame transform uses AES S-boxes and affine 128-bit stages. It
does not execute the generalized lookup-table model. Regenerate its public
profile from the checked research model with:

```sh
UNFAIRPLAY_FRAME_SPN_OUTPUT="$PWD/internal/recovery/profiles/current-frame-spn.bin.gz" \
go test ./internal/recovery -run '^TestWriteFrameSPNProfile$' -count=1
```

Expected size: 17,024 bytes. Expected compressed SHA-256:
`f1cb13c0056a4125387f98a575c21cc6129438580e48c76373e333df614ea963`.
The 51,592-byte decoded profile has SHA-256
`22aabdb275527d924176145147a372bb3c5e4325e5060a345833a2302b1af89c`.

## Normalized SIDB producer SPN

The candidate-frame-to-SIDB-key producer is a normalized semantic function
derived from chosen-state controls. It contains no daemon addresses or
register state. Its private research files are not committed:

```sh
UNFAIRPLAY_SIDB_PRODUCER_PROFILE=/private/path/producer.bin.gz \
UNFAIRPLAY_SIDB_KEY_LAYER_PROFILE=/private/path/key-layer.bin \
UNFAIRPLAY_SIDB_KEY_REPORT_INPUT=/private/path/key-controls.json \
UNFAIRPLAY_SIDB_PRODUCER_SPN_OUTPUT="$PWD/internal/recovery/profiles/current-sidb-producer-spn.bin.gz" \
go test ./internal/recovery -run '^TestSIDBProducerSPNResearch$' -count=1
```

Expected size: 1,929 bytes. Expected compressed SHA-256:
`0b1ce0c897216723802c38dcffd81875cf69afd1d2a164e80cb1391d619733ed`.
The 43,392-byte decoded profile has SHA-256
`3522ed6e1621f8821ce49a179882ac90d046fc62023d77287e174c5333609fc1`.
