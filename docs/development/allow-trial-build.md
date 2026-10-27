# Allow-trial installation build

Official LeviLauncher builds require Microsoft to return a `Full` content license before a retail GDK/MSIXVC package can be installed. The `allowtrialinstall` build tag removes that **Full-only** policy and permits a service-returned `Trial` license to continue through key validation.

This variant does not bypass package encryption, fabricate an entitlement, accept arbitrary keys, or disable integrity checks. Microsoft must still return a usable content key that:

- is bound to the launcher's application-owned device;
- passes key-wrap authentication;
- is 32 bytes after unwrapping; and
- matches the selected package's KeyID.

## Automatic Windows build

From the repository root:

```powershell
pwsh -File .\scripts\build-allow-trial-install.ps1
```

Default output:

```text
bin\LeviLauncher-AllowTrial.exe
```

The script verifies/downloads an isolated Go toolchain, installs the pinned Wails CLI in `%LOCALAPPDATA%\LeviLauncherBuildTools`, rebuilds and tests the native loader, runs tagged Go tests and vet, builds the production frontend, embeds the repository's tracked Windows resources, builds the tagged launcher, and prints its SHA-256. Generated native-loader and `.syso` files are restored/removed after the build so the source tree is not dirtied.

The standalone EXE is unsigned unless you sign it separately with your own code-signing certificate, so Windows may identify it as an unknown publisher.

Requirements that the script checks but does not install:

- 64-bit Windows;
- Node.js/npm;
- xmake; and
- Visual Studio C++ build tools discoverable by xmake.

Options:

```powershell
# Choose another standalone EXE path
pwsh -File .\scripts\build-allow-trial-install.ps1 -OutputPath C:\Builds\LeviLauncher.exe

# Faster repeat build after an already verified build
pwsh -File .\scripts\build-allow-trial-install.ps1 -SkipTests -SkipNativeRebuild
```

## Manual build tag

The policy is selected at compile time:

```powershell
go test -tags allowtrialinstall ./...
go vet -tags allowtrialinstall ./...
go build -tags "production,allowtrialinstall" -trimpath -buildvcs=false -ldflags="-H windowsgui" -o bin\LeviLauncher-AllowTrial.exe .
```

Without `allowtrialinstall`, `nativeinstall.DefaultRequireFullLicense` remains `true`, preserving upstream behavior. The tag changes both the launcher UI installer policy and the native diagnostic CLI's default; the CLI's explicit `-require-full` flag can still override its default.
