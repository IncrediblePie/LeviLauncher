# LeviVanced

LeviVanced is an unofficial custom build of [LeviLauncher](https://github.com/LiteLDev/LeviLauncher) for **Minecraft Bedrock Edition (GDK / UWP)** on Windows.

> [!IMPORTANT]
> **LeviVanced does not unlock Minecraft, grant ownership, bypass package encryption, or fabricate a Microsoft entitlement.** Its only licensing-policy change is allowing a Microsoft-returned `Trial` content license to proceed through installation instead of requiring `Full`. Microsoft must still return a valid device-bound key matching the selected package.

> [!WARNING]
> LeviVanced is an independent community fork. It is **not affiliated with, endorsed by, or supported by LeviLauncher, LeviMC, LiteLDev, Mojang Studios, Microsoft, or Xbox**. Minecraft is a trademark of Microsoft Corporation.

## What changes

- Accepts a valid Microsoft-returned `Trial` or `Full` package key during GDK/MSIXVC installation.
- Checks for launcher updates from [LeviVanced Releases](https://github.com/IncrediblePie/LeviVanced/releases).
- Automatically mirrors each upstream release tag, changelog, and Preview/Latest status.
- Publishes portable and installer builds under the upstream-compatible asset names.

## What remains enforced

- Microsoft account authorization and entitlement response
- Application-device binding
- Content-key unwrap authentication and package KeyID matching
- Package hash-tree, metadata-page, and extraction integrity checks

If Microsoft does not return a usable key for the account, device, and package, installation still fails.

## Downloads

- **LeviVanced releases:** https://github.com/IncrediblePie/LeviVanced/releases
- **Original LeviLauncher releases:** https://github.com/LiteLDev/LeviLauncher/releases

Release executables are automated unsigned community builds. Verify the SHA-256 values included in each release note.

## Upstream documentation

- **User docs:** https://docs.levimc.org/LeviLauncher/
- **中文文档：** https://docs.levimc.org/LeviLauncher/zh-CN/

## Requirements

- Windows 10/11
- A legitimate Microsoft account and access recognized by Microsoft's licensing service; LeviVanced does not provide game access
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

This project is licensed under GPL-3.0-only. The full license text is in `COPYING`, and the bundled third-party components are listed in `THIRD_PARTY_NOTICES`.
