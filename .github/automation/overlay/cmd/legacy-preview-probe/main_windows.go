//go:build windows && amd64 && allowtrialinstall

// SPDX-License-Identifier: GPL-3.0-only
// Diagnostic only: extracts a local package to a NEW directory. It does not
// launch the game, change launcher settings, or overwrite an existing instance.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/liteldev/LeviLauncher/internal/legacyextractor"
	"github.com/liteldev/LeviLauncher/internal/leviloader"
	"github.com/liteldev/LeviLauncher/internal/nativeinstall"
)

func main() {
	input := flag.String("input", "", "Existing local Preview MSIXVC path")
	output := flag.String("output", "", "New output directory (must not exist)")
	timeout := flag.Duration("timeout", 10*time.Minute, "Operation timeout; native call is synchronous")
	flag.Parse()
	if *input == "" || *output == "" {
		flag.Usage()
		os.Exit(2)
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	last := time.Time{}
	err := legacyextractor.Install(ctx, *input, *output, nativeinstall.Options{
		Prepare: func(staging string) error { return leviloader.EnsureForVersion(ctx, staging) },
		Progress: func(p nativeinstall.Progress) {
			if time.Since(last) >= 5*time.Second {
				fmt.Printf("Progress: %d/%d %s\n", p.Current, p.Total, p.File)
				last = time.Now()
			}
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Legacy extraction and current-loader preparation passed: %s\n", *output)
}
