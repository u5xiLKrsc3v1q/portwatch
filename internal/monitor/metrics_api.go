package monitor

import (
	"encoding/json"
	"net/http"
)

// MetricsAPI exposes runtime metrics over HTTP.
type MetricsAPI struct {
	metrics *Metrics
}

// NewMetricsAPI creates a MetricsAPI backed by the given Metrics.
func NewMetricsAPI(m *Metrics) *MetricsAPI {
	return &MetricsAPI{metrics: m}
}

// RegisterRoutes attaches the metrics endpoint to the given mux.
func (a *MetricsAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/metrics", a.handleMetrics)
}

func (a *MetricsAPI) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	snap := a.metrics.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snap); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}
