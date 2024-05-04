package cmd

import (
	"log/slog"

	"github.com/krystal/zadara-exporter/metrics"
	"github.com/krystal/zadara-exporter/zadara/commandcenter"
	"github.com/spf13/cobra"
)

// NewServerCommand creates a new server command for the zadara-exporter.
func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the Zadara exporter server",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := metrics.SetupPrometheusExporter(); err != nil {
				slog.Error("error setting up prometheus exporter", "error", err)

				return
			}

			zClient := commandcenter.NewClientFromToken("baseURL", "token")

			if err := metrics.RegisterStorageMetrics(zClient); err != nil {
				slog.Error("error registering storage metrics", "error", err)

				return
			}

			if err := metrics.Serve(cmd.Context()); err != nil {
				slog.Error("error serving metrics", "error", err)

				return
			}
		},
	}

	return cmd
}
