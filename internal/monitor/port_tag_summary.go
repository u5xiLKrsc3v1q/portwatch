package monitor

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// PortTagSummaryStore tracks counts of entries by tag across scans.
type PortTagSummaryStore struct {
	mu     sync.RWMutex
	counts map[string]int
}

// NewPortTagSummaryStore creates an empty PortTagSummaryStore.
func NewPortTagSummaryStore() *PortTagSummaryStore {
	return &PortTagSummaryStore{
		counts: make(map[string]int),
	}
}

// Record increments the count for each tag present on the given entries.
func (s *PortTagSummaryStore) Record(entries []scanner.Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range entries {
		for _, tag := range e.Tags {
			s.counts[tag]++
		}
	}
}

// Reset clears all tag counts.
func (s *PortTagSummaryStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts = make(map[string]int)
}

// Snapshot returns a copy of the current tag counts.
func (s *PortTagSummaryStore) Snapshot() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]int, len(s.counts))
	for k, v := range s.counts {
		copy[k] = v
	}
	return copy
}

// NewPortTagSummaryAPI returns an http.Handler that serves tag summary data.
func NewPortTagSummaryAPI(store *PortTagSummaryStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Snapshot())
	})
}
