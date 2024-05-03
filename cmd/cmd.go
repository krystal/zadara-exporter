// Package cmd : provides the command line interface for the zadara-exporter.
package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand : creates a new root command for the zadara-exporter.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zadara-exporter",
		Short: "Zadara exporter for Prometheus",
	}

	cmd.AddCommand(NewServerCommand())

	return cmd
}
