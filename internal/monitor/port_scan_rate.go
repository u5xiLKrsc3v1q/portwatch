package monitor

import (
	"net/http"
	"sort"
	"sync"
	"time"
)

// PortScanRateEntry records how frequently a port has been seen across scans.
type PortScanRateEntry struct {
	Port      uint16    `json:"port"`
	Protocol  string    `json:"protocol"`
	ScanCount int       `json:"scan_count"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	RatePerMin float64  `json:"rate_per_min"`
}

// PortScanRateStore tracks how often each port appears across scan cycles.
type PortScanRateStore struct {
	mu      sync.Mutex
	entries map[string]*PortScanRateEntry
}

// NewPortScanRateStore creates a new PortScanRateStore.
func NewPortScanRateStore() *PortScanRateStore {
	return &PortScanRateStore{
		entries: make(map[string]*PortScanRateEntry),
	}
}

func portScanRateKey(port uint16, protocol string) string {
	return protocol + ":" + itoa(int(port))
}

// Record increments the scan count for the given port/protocol pair.
func (s *PortScanRateStore) Record(port uint16, protocol string, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := portScanRateKey(port, protocol)
	e, ok := s.entries[key]
	if !ok {
		s.entries[key] = &PortScanRateEntry{
			Port:      port,
			Protocol:  protocol,
			ScanCount: 1,
			FirstSeen: now,
			LastSeen:  now,
			RatePerMin: 0,
		}
		return
	}
	e.ScanCount++
	e.LastSeen = now
	elapsed := e.LastSeen.Sub(e.FirstSeen).Minutes()
	if elapsed > 0 {
		e.RatePerMin = float64(e.ScanCount) / elapsed
	}
}

// Snapshot returns a sorted copy of all entries.
func (s *PortScanRateStore) Snapshot() []PortScanRateEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]PortScanRateEntry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ScanCount != out[j].ScanCount {
			return out[i].ScanCount > out[j].ScanCount
		}
		return out[i].Port < out[j].Port
	})
	return out
}

// NewPortScanRateAPI returns an HTTP handler for the scan rate store.
func NewPortScanRateAPI(store *PortScanRateStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, store.Snapshot())
	})
}

// itoa is a simple int-to-string helper to avoid importing strconv directly.
func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
