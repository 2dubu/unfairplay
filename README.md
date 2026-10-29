# UnfairPlay

UnfairPlay recovers a FairPlay model from an enrolled Mac and uses it to
decrypt personalized iOS IPAs. It is a standalone Go program.

It does not log in to Apple, buy apps, or download IPAs. You can use ipatool or
another downloader for that part.

## Legal and ethical use

Use UnfairPlay only with apps you lawfully acquired and accounts and devices
you control. It is intended for interoperability, preservation, and authorized
research. The authors do not condone piracy, credential theft, copyright
infringement, or unauthorized distribution of decrypted apps.

## Build

```sh
go build -o unfairplay ./cmd/unfairplay
```

## Usage

First, sign in to the native Mac App Store. Install an iPhone or iPad app and
open it once. Download a matching encrypted IPA with the same Apple Account.

Set up reusable decryption once:

```sh
sudo ./unfairplay setup \
  --app /Applications/App.app \
  --ipa App.ipa
```

Large enrollment databases print bounded recovery progress to standard error.
The final JSON result remains on standard output.

The saved model is private. It is tied to the Mac App Store enrollment, not to
one app, so it can decrypt other personalized IPAs from the same enrollment.

Decrypt another IPA:

```sh
./unfairplay decrypt --input AnotherApp.ipa
```

The default model directory is inside the user configuration directory. Set
`UNFAIRPLAY_MODEL_DIR` or pass `--model` to use another location.

IPAs with extensions and watch apps are decrypted in one pass.

The output keeps the input's `SC_Info` licenses and `iTunesMetadata.plist`,
which identify the Store account. To share the output without them, add
`--strip-account-data`: it omits every `SC_Info` member and removes the Apple
ID and purchaser fields from `iTunesMetadata.plist`. Test without
keeping plaintext with:

```sh
./unfairplay verify --input AnotherApp.ipa
```

The beta also decrypts ARMv7 slices. It validates Thumb-2 or classic ARM
instruction structure and fails closed when a page key produces implausible
code. ARMv6/Thumb-1-only apps are not yet verified.

### Convert an IPA for macOS

This requires an encrypted IPA personalized for the same Store enrollment as
your saved model. Convert it with:

```sh
./unfairplay convert-ipa-macos \
  --input App.ipa \
  --model /private/path/fairplay-enrollment.bin
```

Omit `--model` to use the saved model. The default output is `App.macos.ipa`.
Use `--output` to choose another path or `--force` to replace an existing
output.

The command uses a saved version-2 enrollment model. You can instead use a
matching private 16-byte credential:

```sh
chmod 600 credential.bin
./unfairplay convert-ipa-macos \
  --input App.ipa \
  --credential credential.bin
```

The conversion:

- decodes the current private record with the direct enrollment model;
- validates and rewrites the personalized SINF into the macOS layout;
- preserves its app key, signature, rights, and other private fields;
- adds every SINF copy declared by Apple for encrypted frameworks and
  extensions; and
- leaves all Mach-O executables byte-identical and encrypted.

Model-based conversion supports the observed 440- and 448-byte private SINF
records and the same bounded appended-rights geometry used by recovery.

The result is still an IPA, not a macOS installer. A macOS iOS-app wrapper or
compatible installer must preserve this basic layout:

```text
App iOS.app/
  WrappedBundle -> Wrapper/App.app
  Wrapper/App.app
  Wrapper/iTunesMetadata.plist
```

The converted IPA and wrapper contain a private account license. Keep them
private. This feature does not decrypt, modify, re-sign, or bypass checks inside
the app. Apps may still reject macOS, require unavailable hardware, or depend on
services that are not present on macOS.

## Advanced commands

`collect` and `recover` expose the two setup stages for diagnostics and
archived captures. `doctor` checks a capture, IPA, or model.

```sh
sudo ./unfairplay collect --app /Applications/App.app --ipa App.ipa --output App.capture
./unfairplay recover --capture App.capture --output App.recovery
```

## FAQ

### What is an enrollment?

An enrollment is the FairPlay state created when an Apple Account uses the Mac
App Store on a Mac. It links the account, Mac, and FairPlay selector. `setup`
turns that state into a reusable private model.

### What are SINF, SUPF, SIDB, and SC Info?

- **Personalized SINF:** an app license created for one Store enrollment. It
  identifies the enrollment and protects the app-specific key.
- **SUPF:** per-executable records used to derive keys and IVs for encrypted
  Mach-O pages. Apps, extensions, frameworks, and watch apps may each have one.
- **`SC_Info`:** the directory inside an installed app or IPA that contains its
  SINF and SUPF license files.
- **System `SC Info`:** the private Mac FairPlay state directory. Its SIDB and
  SIDS files contain enrollment state, not app content.

These files are private. Do not publish or share them.

### Do I need to run setup for every app?

No. One model can decrypt different IPAs personalized for the same enrollment.
Each app still has its own key inside its personalized SINF.

### Why might another macOS build fail?

Apple can change enrollment state layouts or semantic tables. UnfairPlay
discovers bounded SIDB and SINF layouts structurally, but fails closed if the
version or authenticated record format is genuinely different. It never
publishes a model until the matching IPA decrypts and passes validation.

Patches and test reports are welcome. Include the macOS build, UnfairPlay
commit, and complete error. Never share SIDB, SINF, hardware
identity, enrollment models, decrypted apps, or Apple credentials.

### Can decryption run on another operating system?

Yes. Collection must start on an enrolled Mac. You can either run `setup` on
that Mac and securely copy the resulting model, or copy a complete private
`collect` bundle and run `recover` elsewhere. The bundle contains the installed
license files, system SC Info state, built-in IOMACAddress, matching IPA, and
macOS build metadata.

The recovered model and personalized IPAs can then be used on macOS or Linux;
fairplayd is never executed. Windows builds compile, but Windows decryption is
not enabled until equivalent private-file permission checks are implemented.

## Support

Recovery is not selected by the macOS build or a fairplayd hash. It uses typed
Go algorithms and does not load or execute fairplayd. Unknown inputs are
accepted only when structural recovery and full IPA validation succeed.

All three observed Mac SIDB geometries are supported: 630, 993, and 1003 bytes.
Other containers in the same semantic family are discovered from their codec,
encrypted registry, record/footer, and trailer geometry. Current SINF records
support the observed short and standard forms plus bounded appended rights.
Every discovered layout still requires complete SINF authentication and IPA
validation.

macOS arm64 and Linux amd64/arm64 are supported. Windows is compile-tested but
cannot yet enforce the required private-file permissions.

The IPA, installed app, Mac state, hardware identity, and model must belong to
the same enrollment. Account or selector mismatches are rejected.

Keep capture and model files private. Do not share SIDB, SINF, hardware data,
recovery models, or decrypted apps. Decrypted IPAs are analysis files and need
signing before installation.

## Branch beta testing

The build-independent beta is on `research/direct-enrollment-algorithm`.
Build it from source and run `setup` normally. Reports should include only the
macOS build, commit, SIDB byte count, SINF `priv`/`righ` byte counts, and the
complete error. Never attach SIDB, SINF, hardware identity, models, credentials,
or decrypted IPAs to a public issue.

## Releases

Run the private control matrix in [QA.md](QA.md) before tagging a release.

Tags matching `v*` publish binaries and `SHA256SUMS` through GitHub Actions.
Tags such as `v0.2.0-rc.1` create prereleases.

## License

MIT. Third-party notices are in `NOTICE`, `LICENSE-DARKPHOENIX.txt`, and
`LICENSE-MACOSFAIRFREE.txt`.
