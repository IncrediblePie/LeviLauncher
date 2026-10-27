//go:build windows && !allowtrialinstall

package nativeinstall

import "testing"

func TestDefaultInstallPolicyRequiresFullLicense(t *testing.T) {
	if !DefaultRequireFullLicense {
		t.Fatal("official build unexpectedly permits trial installation")
	}
}
