// Package cmd provides the command line interface for the zadara-exporter.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRootCommand creates a new root command for the zadara-exporter.
func NewRootCommand(version string, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "zadara-exporter",
		Short:   "Zadara exporter for Prometheus",
		Version: fmt.Sprintf("%s (%s)", version, commit),
	}

	cmd.PersistentFlags().String("config", "", "The path to the configuration file")

	must(viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config")))

	cmd.AddCommand(NewServerCommand())

	return cmd
}
