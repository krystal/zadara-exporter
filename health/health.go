// Package health provides a simple healthcheck handler for HTTP servers.
package health

import (
	"log/slog"
	"net/http"

	"github.com/krystal/zadara-exporter/config"
	"github.com/krystal/zadara-exporter/zadara/commandcenter"
)

// Handler is an HTTP handler for healthchecking purposes.
type Handler struct{}

// ServeHTTP handles the HTTP request and returns a health status.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targets, err := config.GetTargets()
	if err != nil {
		slog.Error("Error getting targets", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var couldNotConnect bool

	for _, target := range targets {
		slog.Info("Checking target", "target", target.Name)
		client := commandcenter.NewClient(target)

		_, err := client.GetStores(r.Context(), target.CloudName)
		if err != nil {
			slog.Error("Error getting stores",
				"name", target.Name,
				"cloud_name", target.CloudName,
				"url", target.URL,
				"error", err)

			couldNotConnect = true
		}
	}

	if couldNotConnect {
		slog.Error("Could not connect to targets", "targets", couldNotConnect)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// RegisterHandler registers the health handler to the provided router.
func RegisterHandler(router *http.ServeMux) {
	handler := &Handler{}
	router.Handle("/healthz", handler)
}
