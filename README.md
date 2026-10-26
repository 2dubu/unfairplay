# UnfairPlay

UnfairPlay recovers a FairPlay model from an enrolled Mac and uses it to
decrypt personalized iOS IPAs. It is a standalone Go program.

It does not log in to Apple, buy apps, or download IPAs. You can use ipatool or
another downloader for that part.

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

The saved model is private. It is tied to the Mac App Store enrollment, not to
one app, so it can decrypt other personalized IPAs from the same enrollment.

Decrypt another IPA:

```sh
./unfairplay decrypt --input AnotherApp.ipa
```

The default model directory is inside the user configuration directory. Set
`UNFAIRPLAY_MODEL_DIR` or pass `--model` to use another location.

IPAs with extensions and watch apps are decrypted in one pass. Test without
keeping plaintext with:

```sh
./unfairplay verify --input AnotherApp.ipa
```

## Advanced commands

`collect` and `recover` expose the two setup stages for diagnostics and
archived captures. `doctor` checks a capture, IPA, or model. `profile` inspects
an arm64e fairplayd without executing it.

```sh
sudo ./unfairplay collect --app /Applications/App.app --ipa App.ipa --output App.capture
./unfairplay recover --capture App.capture --output App.recovery
```

Create a code profile for an exact fairplayd binary:

```sh
./unfairplay profile --binary /path/to/fairplayd.arm64e \
  --recovery-profile-output fairplayd.recovery-profile.gz
```

The profile contains code locations and import names. It contains no account
data and does not execute fairplayd. It helps add native support for new macOS
builds, but creating it alone does not enable decryption.

## Support

Verified enrollment contracts: `25G83`, `26A428`/`26A434`, and `26B118`.
`25G83` and `26B118` use the portable semantic evaluator; native fairplayd
execution is not required. Unknown arm64e builds are accepted only when model
recovery and full IPA validation succeed.

The evaluator source build is recorded for audit purposes. UnfairPlay does not
load or run that fairplayd binary during recovery.

macOS arm64 and Linux amd64/arm64 are supported. Windows is compile-tested but
cannot yet enforce the required private-file permissions.

The IPA, installed app, Mac state, hardware identity, and model must belong to
the same enrollment. Account or selector mismatches are rejected.

Keep capture and model files private. Do not share SIDB, SINF, hardware data,
recovery models, or decrypted apps. Decrypted IPAs are analysis files and need
signing before installation.

## Releases

Run the private control matrix in [QA.md](QA.md) before tagging a release.

Tags matching `v*` publish binaries and `SHA256SUMS` through GitHub Actions.
Tags such as `v0.2.0-rc.1` create prereleases.

## License

MIT. The DarkPhoenix-derived Go code keeps its Apache-2.0 notice in
`LICENSE-DARKPHOENIX.txt`.
