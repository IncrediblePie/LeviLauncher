//go:build windows && amd64 && allowtrialinstall

// SPDX-License-Identifier: GPL-3.0-only
// This wrapper loads the unmodified closed-source upstream v0.3.13 extractor.
package legacyextractor

import (
	"crypto/sha256"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"unicode/utf8"
	"unsafe"

	"golang.org/x/sys/windows"
)

//go:embed launcher_core.dll
var bundledCore []byte

//go:embed NOTICE.txt
var bundledNotice []byte

// CoreSHA256 pins the exact upstream artifact, not a rebuilt approximation.
const CoreSHA256 = "f557c21c08cd89c58bb6d47dd87d27402d4043ad8741b30fc5df91f796ff89aa"

var (
	coreOnce  sync.Once
	coreErr   error
	coreProc  uintptr
	wideToACP = windows.NewLazySystemDLL("kernel32.dll").NewProc("WideCharToMultiByte")
)

func ensureCoreLoaded() error {
	coreOnce.Do(func() {
		cache, err := os.UserCacheDir()
		if err != nil {
			coreErr = err
			return
		}
		path, err := writeCore(filepath.Join(cache, "LeviVanced", "legacy-extractor", CoreSHA256))
		if err != nil {
			coreErr = err
			return
		}
		// Do not accept an environment-provided DLL or search the current directory.
		module, err := windows.LoadLibraryEx(path, 0, windows.LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR|windows.LOAD_LIBRARY_SEARCH_DEFAULT_DIRS)
		if err != nil {
			coreErr = fmt.Errorf("ERR_DLL_DEPENDENCY_MISSING: %w", err)
			return
		}
		coreProc, coreErr = windows.GetProcAddress(module, "GetWithPipe")
		if coreErr != nil {
			_ = windows.FreeLibrary(module)
		}
		// Keep the module loaded for the lifetime of the process.
	})
	return coreErr
}

func writeCore(dir string) (string, error) {
	if fmt.Sprintf("%x", sha256.Sum256(bundledCore)) != CoreSHA256 {
		return "", errors.New("bundled legacy extractor SHA-256 mismatch")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	info, err := os.Lstat(dir)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", errors.New("legacy extractor cache is not a real directory")
	}
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return "", err
	}
	sd, err := windows.SecurityDescriptorFromString("D:P(A;OICI;FA;;;SY)(A;OICI;FA;;;" + user.User.Sid.String() + ")")
	if err != nil {
		return "", err
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		return "", err
	}
	if err := windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, dacl, nil); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, "NOTICE.txt"), bundledNotice, 0600); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "launcher_core.dll")
	if data, err := os.ReadFile(path); err == nil && fmt.Sprintf("%x", sha256.Sum256(data)) == CoreSHA256 {
		return path, nil
	}
	f, err := os.CreateTemp(dir, ".launcher-core-*.dll")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(bundledCore); err != nil {
		_ = f.Close()
		return "", err
	}
	if err = f.Close(); err != nil {
		return "", err
	}
	if err = os.Rename(f.Name(), path); err != nil {
		return "", err
	}
	return path, nil
}

// The old export accepts ANSI strings, not UTF-8. Reject lossy conversions
// instead of silently extracting into a different directory on Unicode paths.
func ansiPath(s string) ([]byte, error) {
	u16, err := windows.UTF16FromString(s)
	if err != nil {
		return nil, err
	}
	if windows.GetACP() == 65001 {
		if !utf8.ValidString(s) {
			return nil, errors.New("invalid UTF-8 path")
		}
		return append([]byte(s), 0), nil
	}
	const noBestFit = 0x400
	var usedDefault uint32
	n, _, callErr := wideToACP.Call(0, noBestFit, uintptr(unsafe.Pointer(&u16[0])), uintptr(len(u16)-1), 0, 0, 0, uintptr(unsafe.Pointer(&usedDefault)))
	if n == 0 {
		return nil, callErr
	}
	b := make([]byte, n+1)
	n, _, callErr = wideToACP.Call(0, noBestFit, uintptr(unsafe.Pointer(&u16[0])), uintptr(len(u16)-1), uintptr(unsafe.Pointer(&b[0])), n, 0, uintptr(unsafe.Pointer(&usedDefault)))
	if n == 0 {
		return nil, callErr
	}
	if usedDefault != 0 {
		return nil, errors.New("legacy extractor cannot represent this path in the Windows ANSI code page")
	}
	return b, nil
}

func callCore(input, output, pipe string) error {
	if err := ensureCoreLoaded(); err != nil {
		return &Error{Code: "ERR_DLL_DEPENDENCY_MISSING", Cause: err}
	}
	args := make([][]byte, 3)
	for i, s := range []string{input, output, pipe} {
		b, err := ansiPath(s)
		if err != nil {
			return &Error{Code: "NH_ERR_INVALID_PARAMS", Cause: err}
		}
		args[i] = b
	}
	rc, _, _ := syscall.SyscallN(coreProc, uintptr(unsafe.Pointer(&args[0][0])), uintptr(unsafe.Pointer(&args[1][0])), uintptr(unsafe.Pointer(&args[2][0])))
	runtime.KeepAlive(args)
	if rc == 0 {
		return nil
	}
	codes := []string{"", "NH_ERR_EXCEPTION", "NH_ERR_INVALID_PARAMS", "NH_ERR_KEY_NOT_FOUND", "NH_ERR_UNAUTHORIZED_CALLER", "NH_ERR_PIPE_OPEN_FAILED", "NH_ERR_INPUT_NOT_FOUND", "NH_ERR_OUTPUT_DIR_INVALID", "NH_ERR_PARSE_FAILED", "NH_ERR_EXTRACT_FAILED"}
	code := "ERR_UNKNOWN"
	if rc < uintptr(len(codes)) {
		code = codes[rc]
	}
	return &Error{Code: code, Cause: fmt.Errorf("upstream legacy extractor returned %d", rc)}
}
