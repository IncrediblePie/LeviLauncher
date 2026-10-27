//go:build windows && !allowtrialinstall

package mcservice

import "context"

func InstallExtractMsixvc(ctx context.Context, name, folderName string, isPreview bool) string {
	return installExtractMsixvcNative(ctx, name, folderName, isPreview)
}
