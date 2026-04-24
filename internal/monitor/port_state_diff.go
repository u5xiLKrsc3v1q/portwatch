package monitor

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/rgst-io/portwatch/internal/scanner"
)

// PortStateDiff represents a single observed change to a port binding.
type PortStateDiff struct {
	Port      int                `json:"port"`
	Protocol  string             `json:"protocol"`
	Address   string             `json:"address"`
	Action    string             `json:"action"` // "added" or "removed"
	ProcessName string           `json:"process_name,omitempty"`
	ObservedAt  time.Time        `json:"observed_at"`
}

// PortStateDiffStore accumulates recent port state diffs for API inspection.
type PortStateDiffStore struct {
	mu      sync.RWMutex
	entries []PortStateDiff
	maxSize int
}

// NewPortStateDiffStore creates a store with the given maximum capacity.
func NewPortStateDiffStore(maxSize int) *PortStateDiffStore {
	if maxSize <= 0 {
		maxSize = 200
	}
	return &PortStateDiffStore{maxSize: maxSize}
}

// Record appends added and removed entries as diffs.
func (s *PortStateDiffStore) Record(added, removed []scanner.Entry) {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range added {
		s.append(PortStateDiff{
			Port:        e.Port,
			Protocol:    e.Protocol.String(),
			Address:     e.Address,
			Action:      "added",
			ProcessName: e.ProcessName,
			ObservedAt:  now,
		})
	}
	for _, e := range removed {
		s.append(PortStateDiff{
			Port:        e.Port,
			Protocol:    e.Protocol.String(),
			Address:     e.Address,
			Action:      "removed",
			ProcessName: e.ProcessName,
			ObservedAt:  now,
		})
	}
}

func (s *PortStateDiffStore) append(d PortStateDiff) {
	if len(s.entries) >= s.maxSize {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, d)
}

// Snapshot returns a copy of all stored diffs.
func (s *PortStateDiffStore) Snapshot() []PortStateDiff {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]PortStateDiff, len(s.entries))
	copy(out, s.entries)
	return out
}

// NewPortStateDiffAPI returns an HTTP handler that serves the diff log.
func NewPortStateDiffAPI(store *PortStateDiffStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(store.Snapshot())
	})
}
