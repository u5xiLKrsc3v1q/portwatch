package monitor

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/rgst-io/portwatch/internal/scanner"
)

// PortSnapshotStore holds the most recent scan snapshot.
type PortSnapshotStore struct {
	mu        sync.RWMutex
	entries   []scanner.Entry
	updatedAt time.Time
}

func NewPortSnapshotStore() *PortSnapshotStore {
	return &PortSnapshotStore{}
}

func (s *PortSnapshotStore) Update(entries []scanner.Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := make([]scanner.Entry, len(entries))
	for i, e := range entries {
		copy[i] = e
	}
	s.entries = copy
	s.updatedAt = time.Now()
}

func (s *PortSnapshotStore) Snapshot() ([]scanner.Entry, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make([]scanner.Entry, len(s.entries))
	for i, e := range s.entries {
		copy[i] = e
	}
	return copy, s.updatedAt
}

// PortSnapshotAPI serves the current port snapshot over HTTP.
type PortSnapshotAPI struct {
	store *PortSnapshotStore
}

func NewPortSnapshotAPI(store *PortSnapshotStore) *PortSnapshotAPI {
	return &PortSnapshotAPI{store: store}
}

func (a *PortSnapshotAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	entries, updatedAt := a.store.Snapshot()
	resp := struct {
		UpdatedAt time.Time       `json:"updated_at"`
		Count     int             `json:"count"`
		Entries   []scanner.Entry `json:"entries"`
	}{
		UpdatedAt: updatedAt,
		Count:     len(entries),
		Entries:   entries,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
