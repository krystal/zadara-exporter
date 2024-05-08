package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/krystal/zadara-exporter/config"
	"github.com/krystal/zadara-exporter/health"
	"github.com/krystal/zadara-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func must(err error) {
	if err != nil {
		slog.Error("error setting up prometheus exporter", "error", err)
		os.Exit(1)
	}
}

func serve(ctx context.Context) error {
	mux := http.NewServeMux()

	// Create a new HTTP handler for serving the metrics.
	mux.Handle(viper.GetString("listen_path"), promhttp.Handler())
	// Register the health handler.
	health.RegisterHandler(mux)

	const ReadHeaderTimeout = 10 * time.Second

	server := &http.Server{
		Addr:              viper.GetString("listen_address"),
		Handler:           mux,
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	slog.Info("starting Metrics server",
		"address", viper.GetString("listen_address"),
		"path", viper.GetString("listen_path"),
	)

	// Start the HTTP server in a separate goroutine.
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("error starting HTTP server: %v", err)
		}
	}()

	// Wait for the context to be done.
	<-ctx.Done()

	// Shutdown the HTTP server gracefully.
	err := server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down HTTP server: %w", err)
	}

	return nil
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

			targets, err := config.GetTargets()
			if err != nil {
				slog.Error("error unmarshalling targets", "error", err)

				return
			}

			if err := metrics.RegisterStorageMetrics(targets); err != nil {
				slog.Error("error registering storage metrics", "error", err)

				return
			}

			if err := serve(cmd.Context()); err != nil {
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
