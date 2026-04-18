package monitor

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// ProcessMapStore tracks the latest known PID/process name per port entry.
type ProcessMapStore struct {
	mu      sync.RWMutex
	entries map[string]scanner.Entry
}

func NewProcessMapStore() *ProcessMapStore {
	return &ProcessMapStore{
		entries: make(map[string]scanner.Entry),
	}
}

func (s *ProcessMapStore) Update(entries []scanner.Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := make(map[string]scanner.Entry, len(entries))
	for _, e := range entries {
		next[e.Key()] = e
	}
	s.entries = next
}

func (s *ProcessMapStore) Snapshot() []scanner.Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]scanner.Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

// ProcessMapAPI exposes current port-to-process mappings over HTTP.
type ProcessMapAPI struct {
	store *ProcessMapStore
}

func NewProcessMapAPI(store *ProcessMapStore) *ProcessMapAPI {
	return &ProcessMapAPI{store: store}
}

func (a *ProcessMapAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	entries := a.store.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
