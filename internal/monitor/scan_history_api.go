package monitor

import (
	"encoding/json"
	"net/http"
)

// ScanHistoryAPI exposes recent scan records over HTTP.
type ScanHistoryAPI struct {
	history *ScanHistory
	mux     *http.ServeMux
}

// NewScanHistoryAPI creates a new ScanHistoryAPI and registers routes on mux.
func NewScanHistoryAPI(history *ScanHistory, mux *http.ServeMux) *ScanHistoryAPI {
	a := &ScanHistoryAPI{history: history, mux: mux}
	mux.HandleFunc("/history", a.handleHistory)
	return a
}

func (a *ScanHistoryAPI) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	records := a.history.Entries()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"count":   len(records),
		"records": records,
	}); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}
