package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Serve : serves the metrics for the zadara-exporter.
func Serve(ctx context.Context) error {
	mux := http.NewServeMux()

	// Create a new HTTP handler for serving the metrics.
	mux.Handle("/metrics", promhttp.Handler())

	const ReadHeaderTimeout = 10 * time.Second

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

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
