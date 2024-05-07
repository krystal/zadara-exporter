package cmd

import (
	"log/slog"
	"os"

	"github.com/krystal/zadara-exporter/config"
	"github.com/krystal/zadara-exporter/metrics"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func must(err error) {
	if err != nil {
		slog.Error("error setting up prometheus exporter", "error", err)
		os.Exit(1)
	}
}

// NewServerCommand creates a new server command for the zadara-exporter.
func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the Zadara exporter server",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := config.Setup(); err != nil {
				slog.Error("error setting up config", "error", err)

				return
			}

			if err := metrics.SetupPrometheusExporter(); err != nil {
				slog.Error("error setting up prometheus exporter", "error", err)

				return
			}

			var targets []*config.Target
			err := viper.UnmarshalKey("targets", &targets)
			if err != nil {
				slog.Error("error unmarshalling targets", "error", err)

				return
			}

			if err := metrics.RegisterStorageMetrics(targets); err != nil {
				slog.Error("error registering storage metrics", "error", err)

				return
			}

			if err := metrics.Serve(cmd.Context()); err != nil {
				slog.Error("error serving metrics", "error", err)

				return
			}
		},
	}

	viper.SetDefault("listen_address", ":9090")
	viper.SetDefault("listen_path", "/metrics")

	cmd.Flags().String("listen_address", ":9090", "The address to listen on for the metrics server")
	cmd.Flags().String("listen_path", "/metrics", "The path to expose the metrics on")

	must(viper.BindPFlag("listen_address", cmd.Flags().Lookup("listen_address")))
	must(viper.BindPFlag("listen_path", cmd.Flags().Lookup("listen_path")))

	return cmd
}
