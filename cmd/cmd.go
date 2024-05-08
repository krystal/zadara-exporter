// Package cmd provides the command line interface for the zadara-exporter.
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var (
// 	Version = "0.0.0-dev"
// ).

// NewRootCommand creates a new root command for the zadara-exporter.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zadara-exporter",
		Short: "Zadara exporter for Prometheus",
	}

	cmd.PersistentFlags().String("config", "", "The path to the configuration file")

	must(viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config")))

	cmd.AddCommand(NewServerCommand())

	return cmd
}
