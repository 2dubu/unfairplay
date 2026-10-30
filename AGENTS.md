# Contributor instructions

These rules apply to the whole repository.

## Safety and privacy

- Never commit Apple credentials, tokens, SIDB/SIDS files, personalized SINF,
  SUPF or SUPP files, enrollment models, hardware identifiers, decrypted apps,
  private captures, or environment-specific absolute paths.
- Tests committed to the repository must use synthetic or public fixtures.
  Private controls belong in the opt-in workflow documented in `QA.md`.
- Keep recovery fail-closed. Do not silently accept an unknown version, schema,
  tag, selector, route, model, page-key class, or malformed signed envelope.
- Do not describe inferred code as verified. State exactly which public,
  synthetic, private, differential, or native controls establish a path.

## Recovery changes

- Production recovery is typed Go. Do not reintroduce native execution,
  fairplayd binaries, lifted CFGs, embedded traces, or build-specific profiles.
- Selection is deterministic: one authenticated container, selected state slot,
  hardware-derived root/tag codebook, exact SINF user/key route, and one value.
  Do not add record, window, key, or license retries to make a fixture pass.
- A new SIDB schema or route needs explicit parsing bounds, fail-closed negative
  tests, complete SINF authentication, and end-to-end IPA validation.
- SUPF/SUPP changes must preserve signed-envelope verification, exact
  architecture/range binding, and plaintext-verification reporting. Missing
  pre-encryption hashes are unverified, never implicitly successful.

## Pull requests

- Keep changes narrow and explain the affected format or enrollment layout.
- Add public or synthetic regression tests for every new branch and rejection.
- Run `gofmt` on changed Go files, then `go test ./...`, `go vet ./...`, and
  `go test -race ./...`.
- When private controls are available, run the matrix from `QA.md` and report
  only anonymous control labels and non-identifying results.
- Review the final diff and repository status. Do not include temporary dumps,
  generated binaries, editor files, or unrelated local changes.
