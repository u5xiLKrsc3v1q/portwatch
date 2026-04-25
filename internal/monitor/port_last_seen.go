package monitor

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PortLastSeenStore tracks the most recent time each port was observed.
type PortLastSeenStore struct {
	mu      sync.RWMutex
	entries map[string]time.Time
}

// NewPortLastSeenStore creates an empty PortLastSeenStore.
func NewPortLastSeenStore() *PortLastSeenStore {
	return &PortLastSeenStore{
		entries: make(map[string]time.Time),
	}
}

// Record updates the last-seen timestamp for each entry.
func (s *PortLastSeenStore) Record(entries []scanner.Entry) {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range entries {
		s.entries[e.Key()] = now
	}
}

// LastSeen returns the last-seen time for the given entry key.
// The second return value is false if the key has never been seen.
func (s *PortLastSeenStore) LastSeen(key string) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.entries[key]
	return t, ok
}

// Snapshot returns a copy of all last-seen timestamps.
func (s *PortLastSeenStore) Snapshot() map[string]time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]time.Time, len(s.entries))
	for k, v := range s.entries {
		copy[k] = v
	}
	return copy
}

// NewPortLastSeenAPI returns an http.Handler that serves last-seen data as JSON.
func NewPortLastSeenAPI(store *PortLastSeenStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		snap := store.Snapshot()
		result := make(map[string]string, len(snap))
		for k, t := range snap {
			result[k] = t.UTC().Format(time.RFC3339)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
}
