// Package cmd provides the command line interface for the zadara-exporter.
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRootCommand creates a new root command for the zadara-exporter.
func NewRootCommand(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "zadara-exporter",
		Short:   "Zadara exporter for Prometheus",
		Version: fmt.Sprintf("%s (%s)", version, commit),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			levelStr := viper.GetString("log-level")

			level := slog.LevelInfo
			err := level.UnmarshalText([]byte(levelStr))
			if err != nil {
				return fmt.Errorf("failed to parse log level: %w", err)
			}

			slog.SetLogLoggerLevel(level)

			return nil
		},
	}

	cmd.PersistentFlags().String("config", "", "The path to the configuration file")
	cmd.PersistentFlags().String("log-level", "info", "The path to the configuration file")

	must(viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config")))
	must(viper.BindPFlag("log-level", cmd.PersistentFlags().Lookup("log-level")))

	cmd.AddCommand(NewServerCommand())

	return cmd
}
