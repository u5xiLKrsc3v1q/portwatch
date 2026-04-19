package monitor

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
)

// PortHitCount tracks how many times a port has been seen across scans.
type PortHitCount struct {
	Port  string `json:"port"`
	Proto string `json:"proto"`
	Count int    `json:"count"`
}

// TopPortsStore accumulates port observation counts.
type TopPortsStore struct {
	mu     sync.Mutex
	counts map[string]*PortHitCount
}

func NewTopPortsStore() *TopPortsStore {
	return &TopPortsStore{
		counts: make(map[string]*PortHitCount),
	}
}

func (s *TopPortsStore) Record(port, proto string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := proto + ":" + port
	if _, ok := s.counts[key]; !ok {
		s.counts[key] = &PortHitCount{Port: port, Proto: proto}
	}
	s.counts[key].Count++
}

func (s *TopPortsStore) Snapshot(limit int) []PortHitCount {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]PortHitCount, 0, len(s.counts))
	for _, v := range s.counts {
		result = append(result, *v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	if limit > 0 && len(result) > limit {
		return result[:limit]
	}
	return result
}

// TopPortsAPI serves the top ports endpoint.
type TopPortsAPI struct {
	store *TopPortsStore
}

func NewTopPortsAPI(store *TopPortsStore) *TopPortsAPI {
	return &TopPortsAPI{store: store}
}

func (a *TopPortsAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a.store.Snapshot(20))
}
