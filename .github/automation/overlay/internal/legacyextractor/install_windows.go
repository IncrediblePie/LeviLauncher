//go:build windows && amd64 && allowtrialinstall

// SPDX-License-Identifier: GPL-3.0-only
package legacyextractor

import (
	"bufio"
	"context"
	"crypto/rand"
	"debug/pe"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Microsoft/go-winio"
	"github.com/beevik/etree"
	"github.com/liteldev/LeviLauncher/internal/nativeinstall"
	"golang.org/x/sys/windows"
)

type Error struct {
	Code  string
	Cause error
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %v", e.Code, e.Cause) }
func (e *Error) Unwrap() error { return e.Cause }

var installGate = make(chan struct{}, 1)

// Install retains the v0.3.13 export but commits only complete, prepared output.
// The opaque synchronous DLL cannot be interrupted safely: cancellation during
// its call waits for it to return, then discards staging rather than committing.
func Install(ctx context.Context, input, output string, options nativeinstall.Options) error {
	return installWithExtractor(ctx, input, output, options, extractWithPipe)
}

type extractFunc func(context.Context, string, string, func(nativeinstall.Progress)) error

func installWithExtractor(ctx context.Context, input, output string, options nativeinstall.Options, extract extractFunc) error {
	if ctx == nil {
		return errors.New("nil context")
	}
	select {
	case installGate <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}
	defer func() { <-installGate }()
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(input) == "" || strings.TrimSpace(output) == "" {
		return &Error{Code: "NH_ERR_INVALID_PARAMS", Cause: errors.New("empty input or output path")}
	}
	in, err := filepath.Abs(input)
	if err != nil {
		return err
	}
	info, err := os.Stat(in)
	if err != nil || !info.Mode().IsRegular() {
		return &Error{Code: "NH_ERR_INPUT_NOT_FOUND", Cause: errors.New("input is not a regular file")}
	}
	out, err := filepath.Abs(output)
	if err != nil || filepath.Dir(out) == out {
		return &Error{Code: "NH_ERR_OUTPUT_DIR_INVALID", Cause: errors.New("invalid output directory")}
	}
	if _, err = os.Lstat(out); !os.IsNotExist(err) {
		return &Error{Code: "ERR_TARGET_EXISTS", Cause: errors.New("output directory already exists or cannot be inspected")}
	}
	if err = os.MkdirAll(filepath.Dir(out), 0700); err != nil {
		return err
	}
	staging, err := os.MkdirTemp(filepath.Dir(out), ".levivanced-preview-legacy-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)
	if err = extract(ctx, in, staging, options.Progress); err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	if err = validateOutput(staging); err != nil {
		return &Error{Code: "NH_ERR_EXTRACT_FAILED", Cause: err}
	}
	if options.Prepare != nil {
		if err = options.Prepare(staging); err != nil {
			return err
		}
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	if _, err = os.Lstat(out); !os.IsNotExist(err) {
		return &Error{Code: "ERR_TARGET_EXISTS", Cause: errors.New("output appeared before commit")}
	}
	return os.Rename(staging, out)
}

func validateOutput(dir string) error {
	for _, name := range []string{"Minecraft.Windows.exe", "MicrosoftGame.config"} {
		info, err := os.Lstat(filepath.Join(dir, name))
		if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
			return fmt.Errorf("missing or invalid extracted %s", name)
		}
	}
	f, err := pe.Open(filepath.Join(dir, "Minecraft.Windows.exe"))
	if err != nil {
		return err
	}
	defer f.Close()
	if f.Machine != pe.IMAGE_FILE_MACHINE_AMD64 {
		return errors.New("extracted game is not Windows x64")
	}
	config, err := os.ReadFile(filepath.Join(dir, "MicrosoftGame.config"))
	if err != nil || len(config) > 1<<20 {
		return errors.New("cannot read extracted game configuration")
	}
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(config); err != nil || doc.Root() == nil || doc.Root().Tag != "Game" {
		return errors.New("invalid extracted game configuration")
	}
	identities := doc.Root().SelectElements("Identity")
	if len(identities) != 1 || identities[0].SelectAttrValue("Name", "") != "Microsoft.MinecraftWindowsBeta" {
		return errors.New("legacy Preview path rejected a non-Preview package")
	}
	return nil
}

func extractWithPipe(ctx context.Context, input, output string, progress func(nativeinstall.Progress)) error {
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return err
	}
	pipe := `\\.\pipe\levivanced_legacy_` + rand.Text()
	listener, err := winio.ListenPipe(pipe, &winio.PipeConfig{SecurityDescriptor: "D:P(A;;GA;;;SY)(A;;GA;;;" + user.User.Sid.String() + ")"})
	if err != nil {
		return &Error{Code: "ERR_CREATE_PIPE", Cause: err}
	}
	var mu sync.Mutex
	var conn net.Conn
	closing := false
	done := make(chan struct{})
	go func() {
		defer close(done)
		c, err := listener.Accept()
		if err != nil {
			return
		}
		mu.Lock()
		if closing {
			mu.Unlock()
			_ = c.Close()
			return
		}
		conn = c
		mu.Unlock()
		defer c.Close()
		scanner := bufio.NewScanner(c)
		scanner.Buffer(make([]byte, 4096), 1<<20)
		for scanner.Scan() {
			if p, ok := decodeProgress(scanner.Bytes()); ok && progress != nil && ctx.Err() == nil {
				progress(p)
			}
		}
	}()
	err = callCore(input, output, pipe)
	_ = listener.Close()
	mu.Lock()
	closing = true
	if conn != nil {
		_ = conn.Close()
	}
	mu.Unlock()
	<-done
	return err
}

func decodeProgress(data []byte) (nativeinstall.Progress, bool) {
	var p struct {
		File          string `json:"file"`
		Current       int64  `json:"current"`
		Total         int64  `json:"total"`
		GlobalCurrent int64  `json:"global_current"`
		GlobalTotal   int64  `json:"global_total"`
	}
	err := json.Unmarshal(data, &p)
	hasControl := func(s string) bool {
		return strings.IndexFunc(s, func(r rune) bool { return r < 32 || r == 127 }) >= 0
	}
	if err != nil || hasControl(p.File) {
		// v0.3.13 also accepted legacy messages with unescaped Windows paths.
		// Some paths accidentally form valid JSON escapes (such as \r and \b);
		// detect the resulting control characters as well as syntax errors.
		if err = json.Unmarshal([]byte(strings.ReplaceAll(string(data), `\`, `/`)), &p); err != nil {
			return nativeinstall.Progress{}, false
		}
	}
	if hasControl(p.File) || p.Current < 0 || p.Total < 0 || p.GlobalCurrent < 0 || p.GlobalTotal < 0 {
		return nativeinstall.Progress{}, false
	}
	return nativeinstall.Progress{File: p.File, Current: p.GlobalCurrent, Total: p.GlobalTotal, FileCurrent: p.Current, FileTotal: p.Total}, true
}
