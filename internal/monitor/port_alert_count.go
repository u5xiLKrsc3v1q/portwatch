package monitor

import (
	"net/http"
	"sort"
	"sync"
	"time"
)

// PortAlertCountStore tracks how many alerts have been fired per port key.
type PortAlertCountStore struct {
	mu      sync.RWMutex
	counts  map[string]int
	lastAt  map[string]time.Time
}

// NewPortAlertCountStore creates an empty PortAlertCountStore.
func NewPortAlertCountStore() *PortAlertCountStore {
	return &PortAlertCountStore{
		counts: make(map[string]int),
		lastAt: make(map[string]time.Time),
	}
}

// Record increments the alert count for the given port key.
func (s *PortAlertCountStore) Record(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[key]++
	s.lastAt[key] = time.Now()
}

// Count returns the total alert count for the given port key.
func (s *PortAlertCountStore) Count(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.counts[key]
}

// PortAlertCountEntry is a snapshot entry for a single port.
type PortAlertCountEntry struct {
	Key    string    `json:"key"`
	Count  int       `json:"count"`
	LastAt time.Time `json:"last_at"`
}

// Snapshot returns a sorted copy of all recorded alert counts.
func (s *PortAlertCountStore) Snapshot() []PortAlertCountEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]PortAlertCountEntry, 0, len(s.counts))
	for k, c := range s.counts {
		out = append(out, PortAlertCountEntry{Key: k, Count: c, LastAt: s.lastAt[k]})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Key < out[j].Key
	})
	return out
}

// Reset clears all alert counts.
func (s *PortAlertCountStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts = make(map[string]int)
	s.lastAt = make(map[string]time.Time)
}

// NewPortAlertCountAPI returns an HTTP handler for the alert count store.
func NewPortAlertCountAPI(store *PortAlertCountStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, store.Snapshot())
	})
}
