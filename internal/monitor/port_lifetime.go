package monitor

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

// PortLifetimeRecord holds the first-seen and last-seen times for a port key.
type PortLifetimeRecord struct {
	Key       string    `json:"key"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	Duration  string    `json:"duration"`
}

// PortLifetimeStore tracks how long each port binding has been observed.
type PortLifetimeStore struct {
	mu      sync.RWMutex
	records map[string]*PortLifetimeRecord
}

// NewPortLifetimeStore creates an empty PortLifetimeStore.
func NewPortLifetimeStore() *PortLifetimeStore {
	return &PortLifetimeStore{
		records: make(map[string]*PortLifetimeRecord),
	}
}

// Record updates the first-seen and last-seen timestamps for a key.
func (s *PortLifetimeStore) Record(key string, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.records[key]; ok {
		r.LastSeen = now
		r.Duration = now.Sub(r.FirstSeen).Truncate(time.Second).String()
	} else {
		s.records[key] = &PortLifetimeRecord{
			Key:       key,
			FirstSeen: now,
			LastSeen:  now,
			Duration:  "0s",
		}
	}
}

// Remove deletes the lifetime record for a key (port gone).
func (s *PortLifetimeStore) Remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.records, key)
}

// Snapshot returns a sorted copy of all lifetime records.
func (s *PortLifetimeStore) Snapshot() []PortLifetimeRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]PortLifetimeRecord, 0, len(s.records))
	for _, r := range s.records {
		out = append(out, *r)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].FirstSeen.Before(out[j].FirstSeen)
	})
	return out
}

// NewPortLifetimeAPI returns an HTTP handler for the lifetime store.
func NewPortLifetimeAPI(store *PortLifetimeStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Snapshot())
	})
}
