# Restored legacy Preview extractor

## Scope

Custom Windows x64 builds (`allowtrialinstall`) use the old `launcher_core.dll` for **Preview GDK/MSIXVC installation only**. Release installation remains on `nativeinstall`, accepting valid Microsoft-returned Trial or Full keys. Builds without the custom tag use the original native installer for both channels. UWP installation, mods/LIP, game runtime licensing, Microsoft account sign-in, and fork launcher-update routing are unchanged.

This restores the exported installation interface used by LeviLauncher v0.3.13. It does **not** grant Minecraft ownership, fabricate entitlements, or patch the game runtime into an unlocked full version. It does mean that the newer WAM content-license request is not used during legacy Preview extraction. The online license badge can therefore show `not_entitled` even when this older extractor is able to install a package.

## Binary provenance

- Upstream: `LiteLDev/LeviLauncher`, release `v0.3.13`.
- Pinned source commit: `5dd6e11ec7694805a96ca83a8379b7ce68a5cac5`.
- Original path: `internal/launchercore/launcher_core.dll`.
- Size: **2,835,968 bytes**.
- SHA-256: `f557c21c08cd89c58bb6d47dd87d27402d4043ad8741b30fc5df91f796ff89aa`.
- Restored export: `GetWithPipe(input ANSI path, output ANSI path, progress pipe ANSI path)`.

[Original caller](https://github.com/LiteLDev/LeviLauncher/blob/5dd6e11ec7694805a96ca83a8379b7ce68a5cac5/internal/launchercore/core.go) and [binary](https://github.com/LiteLDev/LeviLauncher/blob/5dd6e11ec7694805a96ca83a8379b7ce68a5cac5/internal/launchercore/launcher_core.dll).

The binary is copied unchanged, not patched or rebuilt. Its direct imports do not require the old `launcher_api.dll` or `libHttpClient.dll`; those obsolete account/UI dependencies are not restored. The desktop Microsoft Visual C++ runtime is required. The current installer already packages that prerequisite.

## Wrapper behavior

1. Validate the source and destination; reject an existing target.
2. Create a temporary staging directory next to the destination.
3. Verify the embedded DLL's SHA-256, materialize it under a current-user/System-only cache directory, and repair a changed cache copy.
4. Load it by absolute path with restricted Windows DLL-search flags; environment overrides and the current-directory DLL search are not accepted.
5. Pass paths in the current Windows ANSI code page. Reject lossy conversions; UTF-8 system code pages are supported.
6. Read progress from a randomly named, current-user/System-only named pipe.
7. Preserve upstream error codes (including missing key, unauthorized caller, parsing and extraction failures).
8. Require a nonempty x64 `Minecraft.Windows.exe` and a valid `MicrosoftGame.config` identifying `Microsoft.MinecraftWindowsBeta`. Retail packages cannot be committed through the legacy Preview route.
9. Run the current `leviloader.EnsureForVersion` preparation, not v0.3.13's outdated game-loader deployment.
10. Commit the new directory only on success; failures, cancellation, and late target collisions leave no committed partial instance.

The synchronous DLL has no known cancellation export. Cancellation during its call waits for native extraction to return; staging is then discarded. No thread is forcibly terminated in the GUI process.

## Limits and licensing

The old DLL is **closed-source**. Its internal key source, ownership policy, and package hash-tree checks are not established by the public Go caller or by loading its exports. Do not describe it as providing the audited native installer's device-bound entitlement/integrity guarantees. Packages with unsupported formats, unknown keys, or rejected callers can still fail.

The upstream v0.3.13 README states that LGPL-3.0 applies to its **non-closed-source parts**. It does not publish corresponding source for this DLL. The new wrapper is GPL-3.0-only; the unchanged binary's separate rights remain with its upstream authors and it is not relicensed under GPL.

## Automation and verification

The reviewed source/binary overlay lives in `.github/automation/overlay/`. The patcher checks the DLL hash and fails instead of overwriting colliding upstream files. Actions saves the overlay before checking out an exact upstream tag, then applies it alongside the existing custom update-source patch. Portable/installer/updater assets continue using upstream-compatible names.

Tests cover Preview-only routing, native DLL load/export resolution, checksum/cache repair, progress decoding, staging cleanup, preparation failure, cancellation, and preservation of existing/late destinations. A successful build or DLL load is not a successful real Preview installation; a representative package must be tested separately.

### Local verification (2026-10-27)

- Exact upstream `v1.0.4` overlay application and byte-for-byte idempotence passed.
- Legacy wrapper tests and `go vet` passed, including loading the real DLL's `GetWithPipe` export.
- A real local **Preview 1.26.70.20** MSIXVC extracted successfully in approximately 47 seconds and passed current-loader preparation. The resulting configuration identified `Microsoft.MinecraftWindowsBeta`, version `1.26.7020.0`, Store ID `9P5X4QVLC2XR`.
- Extraction used a fresh directory under the approved temporary workspace; the temporary output was removed after verification. Existing launcher instances were not modified. The game itself was not launched, so this is an installation check, not a full-game/runtime-access claim.
- The complete production builder passed native C++ tests, TypeScript/Vite/bundle budgets, the full tagged Go suite, and tagged vet. The full default (non-custom) Go suite and vet also passed.
- Both **Preview 1.26.60.21** and **Preview 1.26.70.20** extracted and passed current-loader preparation with the Preview identity guard enabled (approximately 40 and 48 seconds respectively).
- Inno Setup prerequisite tests, all eight installer languages, wizard layout/cleanup tests, and local installer packaging passed. Portable, installer, and updater ZIP assets were created locally. The final progress-display fix passed regression tests and another real Preview 1.26.70.20 extraction; all artifacts were then rebuilt and repackaged.
- GUI smoke was skipped because an existing launcher was running; it was not stopped or replaced. No Minecraft runtime or full-game access test was performed.
- The inherited frontend dependency audit reported two moderate and two high findings. Frontend dependency versions were not changed or automatically upgraded by this restoration.

### Final local artifacts

Directory: `bin/LeviVanced-LegacyPreview/`. These are unsigned, locally tested artifacts. Automated releases rebuild the same reviewed source in GitHub Actions and publish their own checksums; use the checksum on the release page for downloaded assets.

| Asset | Bytes | SHA-256 |
| --- | ---: | --- |
| `LeviLauncher.exe` | 35,265,024 | `9cf00a2a5ddeeb33ae3b9d62dac2f81d6fbdc1847a4b250a5789698171fd6b6c` |
| `LeviLauncher-amd64-installer.exe` | 39,485,427 | `b90cfb5c4311b823f6848f10c6583d0df2c8980b119b3035552fd2db045bbcf4` |
| `LeviLauncher_windows_amd64.zip` | 20,019,988 | `f9f9e1229f1b599fd41fba0f44d26e9398b7df0ede31e0abac568a51b4f55ce6` |

`SHA256SUMS.txt` and `LEGACY_EXTRACTOR_NOTICE.txt` are alongside these assets. The updater ZIP contains only `LeviLauncher.exe` and its payload matches the portable binary. Keep that executable basename; the existing upstream installer error text also warns against renaming it.

The installer reports version `1.0.4`. Windows FileVersion/ProductVersion fields on the portable executable are blank, as they are on the existing original launcher and prior AllowTrial build; this inherited metadata issue is not introduced by the restoration.

## Local build without committing

The mirror scripts are templates: apply them to an isolated upstream-tag worktree, not to the automation branch's source directory. For example, from the LeviVanced repository in PowerShell:

```powershell
$buildDir = Join-Path $env:LOCALAPPDATA ('Temp\opencode\LeviVanced-local-' + [guid]::NewGuid().ToString('N'))
git worktree add --detach $buildDir v1.0.4
& .\.github\automation\apply-allow-trial-install.ps1 -RepoRoot $buildDir -Prerelease $false
New-Item -ItemType Directory (Join-Path $buildDir 'scripts') -Force | Out-Null
Copy-Item .\.github\automation\build-allow-trial-install.ps1 (Join-Path $buildDir 'scripts\build-allow-trial-install.ps1')
& (Join-Path $buildDir 'scripts\build-allow-trial-install.ps1') -OutputPath 'bin\LeviLauncher.exe'
```

This uses the existing installed build prerequisites, downloads the pinned Go/Wails tools when necessary, and runs the quality gates before building. The executable embeds the legacy DLL; no separate DLL download is required at install time. Keep the asset basename `LeviLauncher.exe`, matching the original exported caller interface and upstream-compatible installer/updater names.

Do not remove the worktree until its build processes have finished. Copy artifacts somewhere permanent before removing a temporary worktree. No commit, push, or public release is needed for this local flow.
