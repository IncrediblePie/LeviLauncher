//go:build windows && allowtrialinstall

// SPDX-License-Identifier: GPL-3.0-only
package nativeinstall

// DefaultRequireFullLicense relaxes the launcher's Full-only policy for builds
// made with -tags allowtrialinstall. The package must still have a valid,
// device-bound content key returned by Microsoft.
const DefaultRequireFullLicense = false
