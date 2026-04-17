package monitor

import (
	"encoding/json"
	"net/http"
	"time"
)

// HTTPAPI exposes a simple read-only HTTP API for portwatch status.
type HTTPAPI struct {
	log      *EventLog
	mux      *http.ServeMux
	server   *http.Server
}

// NewHTTPAPI creates an HTTPAPI bound to addr.
func NewHTTPAPI(addr string, log *EventLog) *HTTPAPI {
	a := &HTTPAPI{log: log, mux: http.NewServeMux()}
	a.server = &http.Server{
		Addr:         addr,
		Handler:      a.mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	a.mux.HandleFunc("/events", a.handleEvents)
	a.mux.HandleFunc("/healthz", a.handleHealth)
	return a
}

// Start begins listening in a goroutine. It returns immediately.
func (a *HTTPAPI) Start() error {
	go func() { _ = a.server.ListenAndServe() }()
	return nil
}

// Stop gracefully shuts down the server.
func (a *HTTPAPI) Stop() error {
	return a.server.Close()
}

func (a *HTTPAPI) handleEvents(w http.ResponseWriter, r *http.Request) {
	events := a.log.Entries()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"count":  len(events),
		"events": events,
	})
}

func (a *HTTPAPI) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
