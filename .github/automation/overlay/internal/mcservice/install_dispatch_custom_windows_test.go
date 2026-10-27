//go:build windows && amd64 && allowtrialinstall

package mcservice

import "testing"

func TestLegacyExtractorOnlyForPreview(t *testing.T) {
	if useLegacyPreviewExtractor(false) {
		t.Fatal("Release installation must retain the native Microsoft-key path")
	}
	if !useLegacyPreviewExtractor(true) {
		t.Fatal("Preview installation must use the restored upstream legacy extractor")
	}
}
