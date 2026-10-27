//go:build windows && allowtrialinstall

package nativeinstall

import "testing"

func TestAllowTrialInstallBuildPolicy(t *testing.T) {
	if DefaultRequireFullLicense {
		t.Fatal("allowtrialinstall build still requires a Full license")
	}
}
