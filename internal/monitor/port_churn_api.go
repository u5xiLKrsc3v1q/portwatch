package monitor

import (
	"encoding/json"
	"net/http"
	"sort"
)

// churnEntry is the JSON shape returned by the churn API.
type churnEntry struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// PortChurnAPI exposes churn scores over HTTP.
type PortChurnAPI struct {
	store *PortChurnStore
}

// NewPortChurnAPI returns an HTTP handler backed by store.
func NewPortChurnAPI(store *PortChurnStore) *PortChurnAPI {
	return &PortChurnAPI{store: store}
}

func (a *PortChurnAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	snap := a.store.Snapshot()
	entries := make([]churnEntry, 0, len(snap))
	for k, v := range snap {
		entries = append(entries, churnEntry{Key: k, Count: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Key < entries[j].Key
	})
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}
