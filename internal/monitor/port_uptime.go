package monitor

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// PortUptimeStore tracks how long each port has been continuously open.
type PortUptimeStore struct {
	mu    sync.RWMutex
	first map[string]time.Time
	now   func() time.Time
}

// NewPortUptimeStore creates a new PortUptimeStore.
func NewPortUptimeStore() *PortUptimeStore {
	return &PortUptimeStore{
		first: make(map[string]time.Time),
		now:   time.Now,
	}
}

// Record marks the first-seen time for a port key if not already recorded.
func (s *PortUptimeStore) Record(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.first[key]; !ok {
		s.first[key] = s.now()
	}
}

// Remove deletes a port key when the port closes.
func (s *PortUptimeStore) Remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.first, key)
}

// Uptime returns how long the given port key has been open.
// Returns 0 and false if the key is unknown.
func (s *PortUptimeStore) Uptime(key string) (time.Duration, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.first[key]
	if !ok {
		return 0, false
	}
	return s.now().Sub(t), true
}

// UptimeRecord is the JSON-serialisable view of a single port's uptime.
type UptimeRecord struct {
	Key     string        `json:"key"`
	Since   time.Time     `json:"since"`
	Uptime  string        `json:"uptime"`
	Seconds float64       `json:"uptime_seconds"`
}

// Snapshot returns all current uptime records.
func (s *PortUptimeStore) Snapshot() []UptimeRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := s.now()
	out := make([]UptimeRecord, 0, len(s.first))
	for k, t := range s.first {
		d := now.Sub(t)
		out = append(out, UptimeRecord{
			Key:     k,
			Since:   t,
			Uptime:  d.Round(time.Second).String(),
			Seconds: d.Seconds(),
		})
	}
	return out
}

// NewPortUptimeAPI returns an http.Handler that serves uptime data as JSON.
func NewPortUptimeAPI(store *PortUptimeStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Snapshot())
	})
}
