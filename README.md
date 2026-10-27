# LeviVanced

LeviVanced is an unofficial custom build of [LeviLauncher](https://github.com/LiteLDev/LeviLauncher) for **Minecraft Bedrock Edition (GDK / UWP)** on Windows.

> [!IMPORTANT]
> **LeviVanced does not grant Minecraft ownership, fabricate Microsoft entitlements, or patch the game runtime to unlock the full game.** Release installations allow valid Microsoft-returned `Trial` keys instead of requiring `Full`; the game's own trial/access restrictions are not removed. Preview installations restore upstream `v0.3.13`'s legacy extractor rather than using the newer WAM content-license gate. This is an installation compatibility change, not a promise of full-game access.

> [!WARNING]
> LeviVanced is an independent community fork. It is **not affiliated with, endorsed by, or supported by LeviLauncher, LeviMC, LiteLDev, Mojang Studios, Microsoft, or Xbox**. Minecraft is a trademark of Microsoft Corporation.

## What changes

- **Release GDK/MSIXVC:** accepts a valid Microsoft-returned `Trial` or `Full` package key.
- **Preview GDK/MSIXVC:** restores the exact upstream `v0.3.13` `launcher_core.dll` extraction path, with progress reporting, transactional output, and current loader preparation. The newer WAM ownership check is not used for this path.
- Checks for launcher updates from [LeviVanced Releases](https://github.com/IncrediblePie/LeviVanced/releases).
- Automatically mirrors each upstream release tag, changelog, and Preview/Latest status.
- Publishes portable and installer builds under the upstream-compatible asset names.

## Validation and limitations

- **Native Release path:** Microsoft account authorization, entitlement response, application-device binding, content-key unwrap authentication, package KeyID matching, and package hash-tree/metadata integrity checks remain enabled.
- **Legacy Preview path:** the DLL is pinned and SHA-256 verified before loading. Output must contain an x64 game executable and game configuration, and loader preparation must succeed before a new instance is committed. Existing instances are never overwritten.
- The legacy DLL is closed-source. Its internal license/key and package-integrity behavior is **not audited**, and must not be represented as equivalent to the native path's checks. New package formats or unknown keys can still fail.
- Cancellation during the synchronous legacy DLL call waits for it to return, then discards the uncommitted output.

Read the [legacy Preview extractor notes](docs/development/legacy-preview-extractor.md) for provenance and scope. The account license badge remains the current Microsoft online/cached report; it does not predict whether the legacy extractor can install a Preview package. UWP installation and launcher self-updates are unchanged.

## Downloads

- **LeviVanced releases:** https://github.com/IncrediblePie/LeviVanced/releases
- **Original LeviLauncher releases:** https://github.com/LiteLDev/LeviLauncher/releases

Release executables are automated unsigned community builds. Verify the SHA-256 values included in each release note.

## Upstream documentation

- **User docs:** https://docs.levimc.org/LeviLauncher/
- **中文文档：** https://docs.levimc.org/LeviLauncher/zh-CN/

## Requirements

- Windows 10/11
- A legitimate licensed copy of Minecraft or an applicable trial; LeviVanced does not provide ownership or remove game-runtime access restrictions
- GDK: Microsoft Gaming Services and Microsoft GameInput available on the system
- UWP: Windows Developer Mode and the UWP framework dependencies required by the package manifest

UWP x64 instances support the console toggle under **Instance settings → Launch options** and compatible `preload-native` DLL mods imported through mod management. Installation deploys the loader; launching a registered instance prepares native loading while preserving Windows app activation and the app channel's shared `LocalState` content. The native loader also requires the desktop Visual C++ runtime, which the launcher checks and offers to install. Vanilla x86/ARM64 packages can still launch, but the bundled native loader currently supports x64 only. Automatic LeviLamina installation remains limited to supported GDK versions.

## Issue Reporting

- **LeviVanced issues:** https://github.com/IncrediblePie/LeviVanced/issues
- Do not request support for LeviVanced from the upstream LeviLauncher/LeviMC communities.
- Include your Windows version, release tag, reproduction steps, and redacted logs/screenshots when possible.

## Development

This repository's end-user guide now lives on GitHub Pages. For local docs work:

```bash
cd docs
npm ci
npm run dev
```

For app development, keep using the existing `frontend/`, `build/`, and Wails task flow. The automated mirror design is documented in [`docs/development/automated-upstream-releases.md`](docs/development/automated-upstream-releases.md).

## License

The current source and new custom wrapper are GPL-3.0-only (`COPYING`). The restored legacy extractor is a separate, unchanged **closed-source upstream binary**, not GPL source; upstream `v0.3.13` limited its LGPL license to non-closed-source parts. The binary retains its upstream authors' separate rights and is not relicensed under GPL. See [legacy extractor provenance](docs/development/legacy-preview-extractor.md) and `THIRD_PARTY_NOTICES`.
