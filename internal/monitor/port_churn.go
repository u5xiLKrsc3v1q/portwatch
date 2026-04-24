package monitor

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PortChurnStore tracks how frequently ports are added and removed within a
// sliding time window, giving a "churn" score per port.
type PortChurnStore struct {
	mu     sync.Mutex
	events map[string][]time.Time
	window time.Duration
}

// NewPortChurnStore returns a PortChurnStore with the given observation window.
func NewPortChurnStore(window time.Duration) *PortChurnStore {
	if window <= 0 {
		window = 5 * time.Minute
	}
	return &PortChurnStore{
		events: make(map[string][]time.Time),
		window: window,
	}
}

// Record registers a churn event for the given entry at the current time.
func (s *PortChurnStore) Record(e scanner.Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := e.Key()
	now := time.Now()
	s.events[key] = append(s.events[key], now)
	s.evict(key, now)
}

// Score returns the number of churn events for the given entry within the window.
func (s *PortChurnStore) Score(e scanner.Entry) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := e.Key()
	s.evict(key, time.Now())
	return len(s.events[key])
}

// Snapshot returns a map of entry key -> churn count for all tracked ports.
func (s *PortChurnStore) Snapshot() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	out := make(map[string]int, len(s.events))
	for k := range s.events {
		s.evict(k, now)
		if len(s.events[k]) > 0 {
			out[k] = len(s.events[k])
		}
	}
	return out
}

// evict removes events outside the observation window. Must be called with lock held.
func (s *PortChurnStore) evict(key string, now time.Time) {
	cutoff := now.Add(-s.window)
	evs := s.events[key]
	i := 0
	for i < len(evs) && evs[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		s.events[key] = evs[i:]
	}
	if len(s.events[key]) == 0 {
		delete(s.events, key)
	}
}
