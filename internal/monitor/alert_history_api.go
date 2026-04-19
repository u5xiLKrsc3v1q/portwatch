package monitor

import (
	"encoding/json"
	"net/http"
)

// AlertHistoryAPI serves alert history over HTTP.
type AlertHistoryAPI struct {
	history *AlertHistory
}

// NewAlertHistoryAPI creates a new AlertHistoryAPI.
func NewAlertHistoryAPI(h *AlertHistory) *AlertHistoryAPI {
	return &AlertHistoryAPI{history: h}
}

func (a *AlertHistoryAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	entries := a.history.Entries()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}
