package monitor

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// PortAgeStore tracks when each port was first seen.
type PortAgeStore struct {
	mu      sync.RWMutex
	firstSeen map[string]time.Time
}

func NewPortAgeStore() *PortAgeStore {
	return &PortAgeStore{
		firstSeen: make(map[string]time.Time),
	}
}

// Record registers a port key if not already seen.
func (s *PortAgeStore) Record(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.firstSeen[key]; !ok {
		s.firstSeen[key] = time.Now()
	}
}

// Age returns how long ago the port was first seen.
func (s *PortAgeStore) Age(key string) (time.Duration, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.firstSeen[key]
	if !ok {
		return 0, false
	}
	return time.Since(t), true
}

// Snapshot returns a copy of all first-seen times.
func (s *PortAgeStore) Snapshot() map[string]time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]time.Time, len(s.firstSeen))
	for k, v := range s.firstSeen {
		out[k] = v
	}
	return out
}

// PortAgeAPI serves port first-seen ages over HTTP.
type PortAgeAPI struct {
	store *PortAgeStore
}

func NewPortAgeAPI(store *PortAgeStore) *PortAgeAPI {
	return &PortAgeAPI{store: store}
}

func (a *PortAgeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	snap := a.store.Snapshot()
	type row struct {
		Key       string    `json:"key"`
		FirstSeen time.Time `json:"first_seen"`
		AgeSeconds float64  `json:"age_seconds"`
	}
	rows := make([]row, 0, len(snap))
	for k, t := range snap {
		rows = append(rows, row{
			Key:        k,
			FirstSeen:  t,
			AgeSeconds: time.Since(t).Seconds(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}
