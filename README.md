# UnfairPlay

UnfairPlay is a standalone, offline Go tool for recovering a native Mac
FairPlay enrollment model and decrypting personalized iOS IPAs for analysis.
It does not log in to Apple, purchase apps, or download IPAs.

The recovery and decryption path does not depend on ipatool, Python, Unicorn,
IDA, a jailbroken iPhone, or an iOS 16 device secret. ipatool is only one
optional way to obtain a personalized encrypted IPA.

## Commands

```text
unfairplay collect   Collect private Mac enrollment state and a matching IPA
unfairplay recover   Recover a model and validate it by decrypting that IPA
unfairplay decrypt   Decrypt another IPA using the recovered model
unfairplay doctor    Diagnose capture/model/IPA compatibility
unfairplay profile   Inspect an arm64e fairplayd without executing it
```

## Build

```sh
go build -o unfairplay ./cmd/unfairplay
```

The resulting executable is self-contained. The pure-Go enrollment evaluator
and its sanitized, account-independent model are embedded in the binary.

## Releases

Pushing a tag matching `v*` runs the full test suite and publishes a GitHub
Release containing versioned binaries for macOS arm64 and Linux amd64/arm64.
Each archive includes the README and license notices. The release also contains
`SHA256SUMS` covering every archive. Windows remains a compile-only CI target
because private-file permission validation is not implemented there.

```sh
git tag -a v0.1.0 -m "UnfairPlay v0.1.0"
git push origin v0.1.0
```

Tags containing a suffix, such as `v0.2.0-rc.1`, are published as prereleases.

## End-to-end workflow

1. Sign in to the native Mac App Store. A fresh Apple Account is the most
   predictable setup, although one preserved older selector-5 enrollment is
   also proven.
2. Install a free iPhone/iPad app through the native Mac App Store and open it
   once. This establishes the Mac enrollment and installed SINF/SUPF files.
3. Obtain a personalized encrypted IPA for the same Apple Account. You may use
   ipatool or another compatible downloader. Purchasing and downloading only
   through ipatool does not replace the native Mac enrollment step.
4. Collect the matching state on the enrolled Mac:

   ```sh
   sudo ./unfairplay collect \
     --app /Applications/App.app \
     --ipa /private/path/App.ipa \
     --output /private/path/App.capture
   ```

5. Recover and validate the enrollment model:

   ```sh
   ./unfairplay recover \
     --capture /private/path/App.capture \
     --output /private/path/App.recovery
   ```

   Recovery publishes `fairplay-enrollment.bin` only after the matching IPA
   passes complete decryption validation. A rejected candidate is not promoted.

6. Reuse that model for another personalized IPA with the same enrollment and
   selector:

   ```sh
   ./unfairplay decrypt \
     --input AnotherApp.ipa \
     --model /private/path/App.recovery/fairplay-enrollment.bin
   ```

## Supported evidence

| Captured FairPlay profile | Status |
| --- | --- |
| `25G83` (`93b8db9c…3da4c5`) | Portable canonical contract verified with an older selector-5 enrollment |
| `26A428` / macOS `26A434` (`1e58d37e…e4ea1784`) | Fully verified pure-Go recovery and IPA decryption |
| `26B118` (`f68e4407…4964439d`) | Verified through the portable 26A428 compatibility evaluator |
| Other arm64e builds using the current SINF/state contract | Attempted, but accepted only after full SINF and matching-IPA validation |

The fairplayd image records provenance and helps classify the capture. Recovery
does not execute that daemon. A nearby build number alone does not prove
compatibility.

## Privacy and limitations

- Capture and recovery directories contain private account-derived state. Keep
  them mode `0700`; model and report files are written mode `0600`.
- Do not share SIDB, hardware identity, SINF, recovery models, or decrypted
  applications.
- The IPA must be personalized for the same Apple Account as the native Mac
  enrollment. Cross-account and cross-machine mixing fails closed.
- Enrollment resets or selector changes require a new capture and model.
- Output IPAs are analysis artifacts. They require appropriate signing before
  installation.
- UnfairPlay does not bypass Apple Account authentication or download access.
  Users remain responsible for lawful use of applications they are authorized
  to access.

## Licensing

The project is MIT licensed; see `LICENSE`. The Go differential-fault recovery
code derived from DarkPhoenix retains its Apache-2.0 notice in
`LICENSE-DARKPHOENIX.txt`.
