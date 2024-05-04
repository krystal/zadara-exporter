package cmd

import (
	"log/slog"
	"os"

	"github.com/krystal/zadara-exporter/config"
	"github.com/krystal/zadara-exporter/metrics"
	"github.com/krystal/zadara-exporter/zadara/commandcenter"
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

			zClient := commandcenter.NewClientFromToken(viper.GetString("api_base_url"), viper.GetString("token"))

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

	viper.SetDefault("api_base_url", "https://api.zadara.com")
	viper.SetDefault("listen_address", ":9090")
	viper.SetDefault("listen_path", "/metrics")

	cmd.Flags().String("api_base_url", "https://api.zadara.com", "The base URL of the Zadara Command Centre API")
	cmd.Flags().String("listen_address", ":9090", "The address to listen on for the metrics server")
	cmd.Flags().String("listen_path", "/metrics", "The path to expose the metrics on")
	cmd.Flags().String("token", "", "The API token for the Zadara Command Centre API")

	must(viper.BindPFlag("api_base_url", cmd.Flags().Lookup("api_base_url")))
	must(viper.BindPFlag("listen_address", cmd.Flags().Lookup("listen_address")))
	must(viper.BindPFlag("listen_path", cmd.Flags().Lookup("listen_path")))
	must(viper.BindPFlag("token", cmd.Flags().Lookup("token")))

	return cmd
}
