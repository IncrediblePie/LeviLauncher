//go:build windows && !allowtrialinstall

// SPDX-License-Identifier: GPL-3.0-only
package nativeinstall

// DefaultRequireFullLicense is the install policy used by the launcher UI and
// diagnostic CLI. Official builds accept only Full content licenses.
const DefaultRequireFullLicense = true
