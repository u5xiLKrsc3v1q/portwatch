package monitor

import (
	"net/http"
)

// RouterConfig holds all API handlers for the HTTP server.
type RouterConfig struct {
	EventLog      *EventLog
	ScanHistory   *ScanHistory
	AlertHistory  *AlertHistory
	PortTrend     *PortTrend
	PortSnapshot  *PortSnapshotStore
	ProcessMap    *ProcessMapStore
	Metrics       *Metrics
}

// NewRouter builds and returns an http.ServeMux wired to all API endpoints.
func NewRouter(cfg RouterConfig) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	if cfg.EventLog != nil {
		mux.Handle("/events", NewHTTPAPI(cfg.EventLog))
	}
	if cfg.ScanHistory != nil {
		mux.Handle("/scans", NewScanHistoryAPI(cfg.ScanHistory))
	}
	if cfg.AlertHistory != nil {
		mux.Handle("/alerts", NewAlertHistoryAPI(cfg.AlertHistory))
	}
	if cfg.PortTrend != nil {
		mux.Handle("/trends", NewPortTrendAPI(cfg.PortTrend))
	}
	if cfg.PortSnapshot != nil {
		mux.Handle("/snapshot", NewPortSnapshotAPI(cfg.PortSnapshot))
	}
	if cfg.ProcessMap != nil {
		mux.Handle("/processes", NewProcessMapAPI(cfg.ProcessMap))
	}
	if cfg.Metrics != nil {
		mux.Handle("/metrics", NewMetricsAPI(cfg.Metrics))
	}

	return mux
}
