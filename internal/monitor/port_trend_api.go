package monitor

import (
	"encoding/json"
	"net/http"
	"sort"
)

// PortTrendAPI exposes port trend data over HTTP.
type PortTrendAPI struct {
	trend *PortTrend
}

// NewPortTrendAPI creates a new PortTrendAPI handler.
func NewPortTrendAPI(trend *PortTrend) *PortTrendAPI {
	return &PortTrendAPI{trend: trend}
}

// ServeHTTP handles GET /trends — returns all port trend entries sorted by SeenCount desc.
func (a *PortTrendAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	entries := a.trend.Snapshot()
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].SeenCount > entries[j].SeenCount
	})
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}
