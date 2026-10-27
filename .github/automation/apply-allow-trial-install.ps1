[CmdletBinding()]
param(
    [string]$RepoRoot = (Get-Location).Path,
    [bool]$Prerelease = $false
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$RepoRoot = [IO.Path]::GetFullPath($RepoRoot)
$utf8 = [Text.UTF8Encoding]::new($false)

function Update-ExactlyOnce {
    param(
        [Parameter(Mandatory)][string]$RelativePath,
        [Parameter(Mandatory)][string]$Before,
        [Parameter(Mandatory)][string]$After
    )

    $path = Join-Path $RepoRoot $RelativePath
    if (-not (Test-Path $path -PathType Leaf)) {
        throw "Required upstream file is missing: $RelativePath"
    }
    $content = [IO.File]::ReadAllText($path)
    $lineEnding = if ($content.Contains("`r`n")) { "`r`n" } else { "`n" }
    $normalizedContent = $content.Replace("`r", "")
    $normalizedBefore = $Before.Replace("`r", "")
    $normalizedAfter = $After.Replace("`r", "")
    $count = [regex]::Matches($normalizedContent, [regex]::Escape($normalizedBefore)).Count
    if ($count -eq 0 -and $normalizedContent.Contains($normalizedAfter)) {
        Write-Host "Already patched: $RelativePath"
        return
    }
    if ($count -ne 1) {
        throw "Expected exactly one patch location in ${RelativePath}; found $count. Review the new upstream release before building."
    }
    $updated = $normalizedContent.Replace($normalizedBefore, $normalizedAfter)
    if ($lineEnding -eq "`r`n") {
        $updated = $updated.Replace("`n", "`r`n")
    }
    [IO.File]::WriteAllText($path, $updated, $utf8)
    Write-Host "Patched: $RelativePath"
}

function Write-PolicyFile {
    param(
        [Parameter(Mandatory)][string]$RelativePath,
        [Parameter(Mandatory)][string]$Content
    )

    $path = Join-Path $RepoRoot $RelativePath
    New-Item -ItemType Directory -Path (Split-Path $path -Parent) -Force | Out-Null
    [IO.File]::WriteAllText($path, $Content.Replace("`n", [Environment]::NewLine), $utf8)
    Write-Host "Wrote: $RelativePath"
}

Update-ExactlyOnce `
    -RelativePath "internal\mcservice\installer.go" `
    -Before 'RequireFullLicense: true,' `
    -After 'RequireFullLicense: nativeinstall.DefaultRequireFullLicense,'

Update-ExactlyOnce `
    -RelativePath "cmd\msixvc-native\main_windows.go" `
    -Before 'flag.Bool("require-full", true, "Reject trial licenses")' `
    -After 'flag.Bool("require-full", nativeinstall.DefaultRequireFullLicense, "Reject trial licenses")'

$innoVersion = @'
#ifndef AppVersion
  #define AppVersion GetVersionNumbersString(AppBinaryPath)
#endif
'@.TrimEnd().Replace("`n", [Environment]::NewLine)
Update-ExactlyOnce `
    -RelativePath "build\windows\inno\project.iss" `
    -Before '#define AppVersion GetVersionNumbersString(AppBinaryPath)' `
    -After $innoVersion

Update-ExactlyOnce `
    -RelativePath "build\windows\inno\project.iss" `
    -Before '#define AppUrl "https://github.com/LiteLDev/LeviLauncher"' `
    -After '#define AppUrl "https://github.com/IncrediblePie/LeviVanced"'

$configPath = Join-Path $RepoRoot "build\config.yml"
$configText = [IO.File]::ReadAllText($configPath)
$betaPattern = '(?m)^(\s*beta:\s*)"(?:true|false)"'
$betaMatches = [regex]::Matches($configText, $betaPattern)
if ($betaMatches.Count -ne 1) {
    throw "Expected exactly one beta setting in build/config.yml; found $($betaMatches.Count)."
}
$betaValue = $Prerelease.ToString().ToLowerInvariant()
$configText = [regex]::Replace($configText, $betaPattern, { param($m) $m.Groups[1].Value + '"' + $betaValue + '"' })
[IO.File]::WriteAllText($configPath, $configText, $utf8)
Write-Host "Set build/config.yml beta to $betaValue."

$releaseAPIFunction = @'
func getReleaseAPIURLs() []string {
	api := "https://api.github.com/repos/" + updateRepository + "/releases"
	return []string{
		api,
		"https://cdn.gh-proxy.org/" + api,
		"https://edgeone.gh-proxy.org/" + api,
		"https://gh-proxy.org/" + api,
		"https://hk.gh-proxy.org/" + api,
		"https://ghproxy.vip/" + api,
	}
}
'@.TrimEnd().Replace("`n", [Environment]::NewLine)
$oldReleaseAPIFunction = @'
func getReleaseAPIURLs() []string {
	return []string{
		"https://api.github.com/repos/LiteLDev/LeviLauncher/releases",
		"https://cdn.gh-proxy.org/https://api.github.com/repos/LiteLDev/LeviLauncher/releases",
		"https://edgeone.gh-proxy.org/https://api.github.com/repos/LiteLDev/LeviLauncher/releases",
		"https://gh-proxy.org/https://api.github.com/repos/LiteLDev/LeviLauncher/releases",
		"https://hk.gh-proxy.org/https://api.github.com/repos/LiteLDev/LeviLauncher/releases",
		"https://ghproxy.vip/https://api.github.com/repos/LiteLDev/LeviLauncher/releases",
	}
}
'@.TrimEnd().Replace("`n", [Environment]::NewLine)
Update-ExactlyOnce -RelativePath "internal\update\update.go" -Before $oldReleaseAPIFunction -After $releaseAPIFunction

$assetFields = @'
Size               int64  `json:"size"`
	Digest             string `json:"digest"`
'@.TrimEnd().Replace("`n", [Environment]::NewLine)
Update-ExactlyOnce `
    -RelativePath "internal\update\update.go" `
    -Before 'Size               int64  `json:"size"`' `
    -After $assetFields

$oldUpdateCandidates = @'
func buildSelfUpdateCandidateURLs(ver string) []string {
	return []string{
		"https://goproxy.io/github.com/dreamguxiang/!levi!launcher/@v/" + ver + ".zip",
		"https://mirrors.aliyun.com/goproxy/github.com/dreamguxiang/levilauncher/@v/" + ver + ".zip",
	}
}
'@.TrimEnd().Replace("`n", [Environment]::NewLine)
$newUpdateCandidates = @'
func buildSelfUpdateCandidateURLs(release githubRelease, ver string) ([]string, string, error) {
	if useReleaseUpdateAsset {
		asset := findAssetByName(release, expectedZipAssetName())
		if asset == nil || strings.TrimSpace(asset.BrowserDownloadURL) == "" {
			return nil, "", fmt.Errorf("release is missing %s", expectedZipAssetName())
		}
		digest := strings.TrimSpace(asset.Digest)
		if !strings.HasPrefix(strings.ToLower(digest), "sha256:") {
			return nil, "", fmt.Errorf("release asset has no SHA-256 digest")
		}
		hash, err := parseSHA256Checksum(strings.TrimSpace(digest[len("sha256:"):]))
		if err != nil {
			return nil, "", fmt.Errorf("invalid release asset digest: %w", err)
		}
		return buildAssetCandidateURLs(asset.BrowserDownloadURL), hash, nil
	}
	return []string{
		"https://goproxy.io/github.com/dreamguxiang/!levi!launcher/@v/" + ver + ".zip",
		"https://mirrors.aliyun.com/goproxy/github.com/dreamguxiang/levilauncher/@v/" + ver + ".zip",
	}, "", nil
}
'@.TrimEnd().Replace("`n", [Environment]::NewLine)
Update-ExactlyOnce -RelativePath "internal\update\update.go" -Before $oldUpdateCandidates -After $newUpdateCandidates

$oldDownloadBlock = @'
	if fi, statErr := os.Stat(zipPath); statErr == nil && fi.Size() > 0 {
		application.Get().Event.Emit(EventAppUpdateStatus, "downloaded")
		application.Get().Event.Emit(EventAppUpdateProgress, AppUpdateProgress{Phase: "download", Downloaded: fi.Size(), Total: fi.Size()})
	} else {
		if err := downloadUpdateAsset(ctx, buildSelfUpdateCandidateURLs(ver), zipPath); err != nil {
			return emitUpdateError("ERR_UPDATE_DOWNLOAD_FAILED", err)
		}
	}
'@.TrimEnd().Replace("`n", [Environment]::NewLine)
$newDownloadBlock = @'
	updateURLs, expectedHash, err := buildSelfUpdateCandidateURLs(release, ver)
	if err != nil {
		return emitUpdateError("ERR_UPDATE_RELEASE_UNAVAILABLE", err)
	}
	cacheValid := false
	if fi, statErr := os.Stat(zipPath); statErr == nil && fi.Size() > 0 {
		cacheValid = expectedHash == "" || verifyFileSHA256(zipPath, expectedHash) == nil
		if cacheValid {
			application.Get().Event.Emit(EventAppUpdateStatus, "downloaded")
			application.Get().Event.Emit(EventAppUpdateProgress, AppUpdateProgress{Phase: "download", Downloaded: fi.Size(), Total: fi.Size()})
		} else {
			_ = os.Remove(zipPath)
		}
	}
	if !cacheValid {
		if err := downloadUpdateAsset(ctx, updateURLs, zipPath); err != nil {
			return emitUpdateError("ERR_UPDATE_DOWNLOAD_FAILED", err)
		}
		if expectedHash != "" {
			if err := verifyFileSHA256(zipPath, expectedHash); err != nil {
				_ = os.Remove(zipPath)
				return emitUpdateError("ERR_UPDATE_DOWNLOAD_FAILED", err)
			}
		}
	}
'@.TrimEnd().Replace("`n", [Environment]::NewLine)
Update-ExactlyOnce -RelativePath "internal\update\update.go" -Before $oldDownloadBlock -After $newDownloadBlock

Write-PolicyFile "internal\nativeinstall\install_policy_default_windows.go" @'
//go:build windows && !allowtrialinstall

// SPDX-License-Identifier: GPL-3.0-only
package nativeinstall

// DefaultRequireFullLicense is the install policy used by the launcher UI and
// diagnostic CLI. Official builds accept only Full content licenses.
const DefaultRequireFullLicense = true
'@

Write-PolicyFile "internal\nativeinstall\install_policy_allow_trial_windows.go" @'
//go:build windows && allowtrialinstall

// SPDX-License-Identifier: GPL-3.0-only
package nativeinstall

// DefaultRequireFullLicense relaxes the launcher's Full-only policy for builds
// made with -tags allowtrialinstall. The package must still have a valid,
// device-bound content key returned by Microsoft.
const DefaultRequireFullLicense = false
'@

Write-PolicyFile "internal\nativeinstall\install_policy_default_windows_test.go" @'
//go:build windows && !allowtrialinstall

package nativeinstall

import "testing"

func TestDefaultInstallPolicyRequiresFullLicense(t *testing.T) {
	if !DefaultRequireFullLicense {
		t.Fatal("official build unexpectedly permits trial installation")
	}
}
'@

Write-PolicyFile "internal\nativeinstall\install_policy_allow_trial_windows_test.go" @'
//go:build windows && allowtrialinstall

package nativeinstall

import "testing"

func TestAllowTrialInstallBuildPolicy(t *testing.T) {
	if DefaultRequireFullLicense {
		t.Fatal("allowtrialinstall build still requires a Full license")
	}
}
'@

Write-PolicyFile "internal\update\source_default.go" @'
//go:build !allowtrialinstall

package update

const (
	updateRepository     = "LiteLDev/LeviLauncher"
	useReleaseUpdateAsset = false
)
'@

Write-PolicyFile "internal\update\source_allow_trial.go" @'
//go:build allowtrialinstall

package update

const (
	updateRepository     = "IncrediblePie/LeviVanced"
	useReleaseUpdateAsset = true
)
'@

Write-PolicyFile "internal\update\source_default_test.go" @'
//go:build !allowtrialinstall

package update

import "testing"

func TestDefaultUpdateSource(t *testing.T) {
	if updateRepository != "LiteLDev/LeviLauncher" || useReleaseUpdateAsset {
		t.Fatal("official build update source changed unexpectedly")
	}
	urls, hash, err := buildSelfUpdateCandidateURLs(githubRelease{}, "v1.2.3")
	if err != nil || hash != "" || len(urls) != 2 {
		t.Fatalf("unexpected official update candidates: urls=%v hash=%q err=%v", urls, hash, err)
	}
}
'@

Write-PolicyFile "internal\update\source_allow_trial_test.go" @'
//go:build allowtrialinstall

package update

import (
	"strings"
	"testing"
)

func TestAllowTrialUpdateSource(t *testing.T) {
	if updateRepository != "IncrediblePie/LeviVanced" || !useReleaseUpdateAsset {
		t.Fatal("allowtrialinstall build is not using the fork release source")
	}
	digest := strings.Repeat("a", 64)
	target := "https://github.com/IncrediblePie/LeviVanced/releases/download/v1.2.3/LeviLauncher_windows_amd64.zip"
	release := githubRelease{Assets: []githubReleaseAsset{{
		Name: expectedZipAssetName(), BrowserDownloadURL: target, Digest: "sha256:" + digest,
	}}}
	urls, hash, err := buildSelfUpdateCandidateURLs(release, "v1.2.3")
	if err != nil || hash != digest {
		t.Fatalf("unexpected custom update candidate: hash=%q err=%v", hash, err)
	}
	found := false
	for _, candidate := range urls {
		found = found || candidate == target
	}
	if !found {
		t.Fatalf("fork release URL missing from candidates: %v", urls)
	}
	if _, _, err := buildSelfUpdateCandidateURLs(githubRelease{}, "v1.2.3"); err == nil {
		t.Fatal("custom updater accepted a release without its update ZIP")
	}
}
'@

Write-Host "Allow-trial installation patch applied successfully." -ForegroundColor Green
