package monitor

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// PortConnectionCountStore tracks how many times each port has been seen active.
type PortConnectionCountStore struct {
	mu     sync.RWMutex
	counts map[string]int
}

// NewPortConnectionCountStore creates a new PortConnectionCountStore.
func NewPortConnectionCountStore() *PortConnectionCountStore {
	return &PortConnectionCountStore{
		counts: make(map[string]int),
	}
}

// Record increments the connection count for each entry in the scan.
func (s *PortConnectionCountStore) Record(entries []scanner.Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range entries {
		s.counts[e.Key()]++
	}
}

// Snapshot returns a copy of the current counts.
func (s *PortConnectionCountStore) Snapshot() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]int, len(s.counts))
	for k, v := range s.counts {
		out[k] = v
	}
	return out
}

// Reset clears all counts.
func (s *PortConnectionCountStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts = make(map[string]int)
}

// PortConnectionCountEntry is the API response shape.
type PortConnectionCountEntry struct {
	Port  string `json:"port"`
	Count int    `json:"count"`
}

// NewPortConnectionCountAPI returns an HTTP handler for connection counts.
func NewPortConnectionCountAPI(store *PortConnectionCountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		snap := store.Snapshot()
		result := make([]PortConnectionCountEntry, 0, len(snap))
		for port, count := range snap {
			result = append(result, PortConnectionCountEntry{Port: port, Count: count})
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].Count > result[j].Count
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
