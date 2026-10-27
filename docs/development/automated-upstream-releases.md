# Automated upstream release mirror

The fork's `Mirror upstream release` workflow checks `LiteLDev/LeviLauncher` every six hours. A scheduled run uses the newest non-draft release only after both upstream Windows assets are present. It exits before allocating a Windows builder if the fork already has that tag with both mirrored asset names.

For each unmirrored release it:

1. mirrors upstream `main` to the fork's dedicated `upstream-main` source-snapshot branch without overwriting LeviVanced's automation branch (upstream workflow files are omitted because GitHub forbids its Actions token from rewriting workflows);
2. checks out the exact upstream release tag in detached mode;
3. applies the reviewed `allowtrialinstall` policy patch with exact-match guards;
4. builds and tests `LeviLauncher.exe`;
5. packages that exact executable as `LeviLauncher-amd64-installer.exe` using the upstream Inno Setup project and as `LeviLauncher_windows_amd64.zip` for in-app updates;
6. copies the upstream release name and changelog;
7. copies the selected upstream release tag only after its assets exist, then publishes a `LeviVanced` release with the same tag and prerelease status; stable releases become the fork's latest release; and
8. uploads the portable and installer assets under the same names used upstream, plus the updater ZIP.

Manual runs accept an optional upstream tag and a force-rebuild switch. Force rebuild replaces the two assets and refreshes release metadata.

Builds made with `allowtrialinstall` query `IncrediblePie/LeviVanced` for updates. Preview builds automatically accept mirrored preview updates; stable builds retain the normal stable-only default. The updater downloads `LeviLauncher_windows_amd64.zip` from the selected fork release and verifies GitHub's SHA-256 asset digest before installation. The installer metadata also points to the fork's Releases page. Ordinary builds without the custom tag continue using the upstream update source.

The release note identifies the binaries as unsigned custom builds and includes their SHA-256 values. Package encryption, Microsoft entitlement responses, device binding, key-wrap authentication, package integrity checks, and KeyID matching remain enabled.

GitHub only evaluates scheduled workflows from a repository's default branch. As with all public GitHub repositories, GitHub may disable schedules after a prolonged period without repository activity; a manual workflow run re-enables/checks the mirror immediately.

The encrypted `LEVIVANCED_WORKFLOW_TOKEN` repository secret is used only while copying the exact annotated upstream release tag. GitHub's built-in Actions token cannot introduce commits containing workflow changes. The workflow clears this token before checking out or executing any upstream release code.
