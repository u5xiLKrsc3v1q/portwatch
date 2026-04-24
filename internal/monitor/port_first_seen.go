package monitor

import (
	"net/http"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PortFirstSeenStore tracks when each port was first observed.
type PortFirstSeenStore struct {
	mu      sync.RWMutex
	records map[string]time.Time
}

// NewPortFirstSeenStore creates an empty PortFirstSeenStore.
func NewPortFirstSeenStore() *PortFirstSeenStore {
	return &PortFirstSeenStore{
		records: make(map[string]time.Time),
	}
}

// Record registers the first-seen time for an entry if not already known.
func (s *PortFirstSeenStore) Record(e scanner.Entry) {
	key := e.Key()
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.records[key]; !ok {
		s.records[key] = time.Now()
	}
}

// FirstSeen returns the time an entry was first seen, and whether it exists.
func (s *PortFirstSeenStore) FirstSeen(e scanner.Entry) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.records[e.Key()]
	return t, ok
}

// Snapshot returns a copy of all first-seen records.
func (s *PortFirstSeenStore) Snapshot() map[string]time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]time.Time, len(s.records))
	for k, v := range s.records {
		out[k] = v
	}
	return out
}

// NewPortFirstSeenAPI returns an HTTP handler that exposes the first-seen store.
func NewPortFirstSeenAPI(store *PortFirstSeenStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		snap := store.Snapshot()
		type row struct {
			Key       string    `json:"key"`
			FirstSeen time.Time `json:"first_seen"`
		}
		rows := make([]row, 0, len(snap))
		for k, t := range snap {
			rows = append(rows, row{Key: k, FirstSeen: t})
		}
		writeJSON(w, rows)
	})
}
