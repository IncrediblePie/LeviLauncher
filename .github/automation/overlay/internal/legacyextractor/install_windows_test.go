//go:build windows && amd64 && allowtrialinstall

package legacyextractor

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/liteldev/LeviLauncher/internal/nativeinstall"
)

func testSource(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	in := filepath.Join(root, "input.msixvc")
	if err := os.WriteFile(in, []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	return in, filepath.Join(root, "installed")
}

func fixtureExtract(t *testing.T) extractFunc {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	return func(_ context.Context, _, staging string, _ func(nativeinstall.Progress)) error {
		if err := os.WriteFile(filepath.Join(staging, "Minecraft.Windows.exe"), data, 0600); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(staging, "MicrosoftGame.config"), []byte(`<Game><Identity Name="Microsoft.MinecraftWindowsBeta"/></Game>`), 0600)
	}
}

func assertNoStaging(t *testing.T, output string) {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(filepath.Dir(output), ".levivanced-preview-legacy-*"))
	if err != nil || len(paths) != 0 {
		t.Fatalf("leftover staging: %v (%v)", paths, err)
	}
}

func TestInstallPreparesBeforeCommit(t *testing.T) {
	in, out := testSource(t)
	prepared := false
	err := installWithExtractor(context.Background(), in, out, nativeinstall.Options{Prepare: func(staging string) error {
		if _, err := os.Stat(out); !os.IsNotExist(err) {
			t.Fatal("output committed before preparation")
		}
		prepared = true
		return os.WriteFile(filepath.Join(staging, "loader-prepared"), []byte("ok"), 0600)
	}}, fixtureExtract(t))
	if err != nil || !prepared {
		t.Fatalf("install: %v; prepared=%v", err, prepared)
	}
	if _, err := os.Stat(filepath.Join(out, "loader-prepared")); err != nil {
		t.Fatal(err)
	}
	assertNoStaging(t, out)
}

func TestInstallPreservesExistingOutput(t *testing.T) {
	in, out := testSource(t)
	if err := os.Mkdir(out, 0700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(out, "keep")
	if err := os.WriteFile(marker, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	called := false
	err := installWithExtractor(context.Background(), in, out, nativeinstall.Options{}, func(context.Context, string, string, func(nativeinstall.Progress)) error {
		called = true
		return nil
	})
	var failure *Error
	if !errors.As(err, &failure) || failure.Code != "ERR_TARGET_EXISTS" || called {
		t.Fatalf("existing target was not rejected: %v; called=%v", err, called)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "keep" {
		t.Fatal("existing output changed")
	}
}

func TestFailedPreparationDiscardsStaging(t *testing.T) {
	in, out := testSource(t)
	want := errors.New("preparation failed")
	err := installWithExtractor(context.Background(), in, out, nativeinstall.Options{Prepare: func(string) error { return want }}, fixtureExtract(t))
	if !errors.Is(err, want) {
		t.Fatalf("got %v", err)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatal("failed install committed output")
	}
	assertNoStaging(t, out)
}

func TestCanceledExtractionDiscardsStaging(t *testing.T) {
	in, out := testSource(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fixture := fixtureExtract(t)
	err := installWithExtractor(ctx, in, out, nativeinstall.Options{}, func(ctx context.Context, in, staging string, progress func(nativeinstall.Progress)) error {
		if err := fixture(ctx, in, staging, progress); err != nil {
			return err
		}
		cancel()
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatal("canceled extraction committed output")
	}
	assertNoStaging(t, out)
}

func TestExtractorFailureAndIncompleteOutput(t *testing.T) {
	for _, mode := range []string{"error", "incomplete"} {
		t.Run(mode, func(t *testing.T) {
			in, out := testSource(t)
			err := installWithExtractor(context.Background(), in, out, nativeinstall.Options{}, func(_ context.Context, _, staging string, _ func(nativeinstall.Progress)) error {
				if err := os.WriteFile(filepath.Join(staging, "partial"), []byte("partial"), 0600); err != nil {
					return err
				}
				if mode == "error" {
					return &Error{Code: "NH_ERR_KEY_NOT_FOUND", Cause: errors.New("fixture")}
				}
				return nil
			})
			if err == nil {
				t.Fatal("failed/incomplete extraction accepted")
			}
			if _, err := os.Stat(out); !os.IsNotExist(err) {
				t.Fatal("failed install committed output")
			}
			assertNoStaging(t, out)
		})
	}
}

func TestLateOutputIsPreserved(t *testing.T) {
	in, out := testSource(t)
	err := installWithExtractor(context.Background(), in, out, nativeinstall.Options{Prepare: func(string) error { return os.Mkdir(out, 0700) }}, fixtureExtract(t))
	var failure *Error
	if !errors.As(err, &failure) || failure.Code != "ERR_TARGET_EXISTS" {
		t.Fatalf("late output not rejected: %v", err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal("late output removed")
	}
	assertNoStaging(t, out)
}

func TestLegacyProgress(t *testing.T) {
	for _, text := range []string{
		`{"file":"C:\\game\\file","current":3,"total":4,"global_current":5,"global_total":6}`,
		`{"file":"C:\game\file","current":3,"total":4,"global_current":5,"global_total":6}`,
	} {
		p, ok := decodeProgress([]byte(text))
		if !ok || p.FileCurrent != 3 || p.FileTotal != 4 || p.Current != 5 || p.Total != 6 {
			t.Fatalf("bad progress: %+v, %v", p, ok)
		}
	}
	for _, text := range []string{`broken`, `{"current":-1}`, strings.Repeat("x", 1024)} {
		if _, ok := decodeProgress([]byte(text)); ok {
			t.Fatal("malformed progress accepted")
		}
	}
}

func TestLegacyPathRejectsRetailPackage(t *testing.T) {
	in, out := testSource(t)
	fixture := fixtureExtract(t)
	err := installWithExtractor(context.Background(), in, out, nativeinstall.Options{}, func(ctx context.Context, in, staging string, progress func(nativeinstall.Progress)) error {
		if err := fixture(ctx, in, staging, progress); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(staging, "MicrosoftGame.config"), []byte(`<Game><Identity Name="Microsoft.MinecraftUWP"/></Game>`), 0600)
	})
	if err == nil {
		t.Fatal("legacy Preview installer accepted a retail package")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatal("retail package committed through legacy Preview path")
	}
	assertNoStaging(t, out)
}

func TestLegacyProgressRepairsAccidentalJSONEscapes(t *testing.T) {
	p, ok := decodeProgress([]byte(`{"file":"data\ray_tracing\blue_noise","current":1,"total":2,"global_current":3,"global_total":4}`))
	if !ok || p.File != "data/ray_tracing/blue_noise" {
		t.Fatalf("legacy path became control characters: %q (ok=%v)", p.File, ok)
	}
	p, ok = decodeProgress([]byte(`{"file":"\u4f60\u597d/file","current":1}`))
	if !ok || p.File != "你好/file" {
		t.Fatalf("valid Unicode JSON was changed: %q (ok=%v)", p.File, ok)
	}
}

func TestCanceledOrNilContextDoesNotExtract(t *testing.T) {
	in, out := testSource(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, ctx := range []context.Context{nil, ctx} {
		called := false
		err := installWithExtractor(ctx, in, out, nativeinstall.Options{}, func(context.Context, string, string, func(nativeinstall.Progress)) error {
			called = true
			return nil
		})
		if err == nil || called {
			t.Fatalf("invalid/canceled context extracted: %v, called=%v", err, called)
		}
	}
}

func TestANSIPathRejectsNUL(t *testing.T) {
	if _, err := ansiPath("invalid\x00path"); err == nil {
		t.Fatal("embedded NUL path accepted")
	}
}

func TestBundledCoreIntegrityAndCacheRepair(t *testing.T) {
	dir := t.TempDir()
	path, err := writeCore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("tampered"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := writeCore(dir); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || fmt.Sprintf("%x", sha256.Sum256(data)) != CoreSHA256 {
		t.Fatal("cache repair did not restore pinned DLL")
	}
}

func TestBundledCoreLoadsAndExportsGetWithPipe(t *testing.T) {
	if err := ensureCoreLoaded(); err != nil {
		t.Fatal(err)
	}
	if coreProc == 0 {
		t.Fatal("GetWithPipe export was not resolved")
	}
}
