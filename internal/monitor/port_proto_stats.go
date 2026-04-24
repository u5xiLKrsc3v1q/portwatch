package monitor

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/yourorg/portwatch/internal/scanner"
)

// ProtoStats holds counts of observed entries per protocol.
type ProtoStats struct {
	TCP uint64 `json:"tcp"`
	UDP uint64 `json:"udp"`
	Other uint64 `json:"other"`
}

// PortProtoStatsStore accumulates per-protocol port counts.
type PortProtoStatsStore struct {
	mu    sync.RWMutex
	stats ProtoStats
}

// NewPortProtoStatsStore returns an initialised store.
func NewPortProtoStatsStore() *PortProtoStatsStore {
	return &PortProtoStatsStore{}
}

// Record increments the counter for each entry's protocol.
func (s *PortProtoStatsStore) Record(entries []scanner.Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range entries {
		switch {
		case e.Protocol.IsTCP():
			s.stats.TCP++
		case e.Protocol.IsUDP():
			s.stats.UDP++
		default:
			s.stats.Other++
		}
	}
}

// Snapshot returns a copy of the current stats.
func (s *PortProtoStatsStore) Snapshot() ProtoStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stats
}

// Reset zeroes all counters.
func (s *PortProtoStatsStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats = ProtoStats{}
}

// NewPortProtoStatsAPI returns an http.Handler that serves protocol stats.
func NewPortProtoStatsAPI(store *PortProtoStatsStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(store.Snapshot())
	})
}
