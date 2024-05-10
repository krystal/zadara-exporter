// Package main provides the entry point for the zadara-exporter.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/krystal/zadara-exporter/cmd"
)

//nolint:gochecknoglobals // These variables are set using ldflags.
var (
	version = "0.0.0-dev"
	commit  = ""
)

func mainE() error {
	rootCmd := cmd.NewRootCommand(version, commit)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	err := rootCmd.ExecuteContext(ctx)

	return fmt.Errorf("command failed: %w", err)
}

func main() {
	if err := mainE(); err != nil {
		os.Exit(1)
	}
}
