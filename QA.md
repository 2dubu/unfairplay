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
      "expected_evaluator_source_sha256": "...",
      "expected_reference_code_source": "embedded-profile",
      "expected_images": 10,
      "expected_pages": 11004
    }
  ]
}
```

The command fails on any changed hash, image count, page count, or incomplete
decryption. Use it before creating a release tag.
