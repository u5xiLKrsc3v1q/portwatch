package monitor

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// AnomalyRecord describes a port that appeared outside of expected hours.
type AnomalyRecord struct {
	Port      uint16    `json:"port"`
	Address   string    `json:"address"`
	Protocol  string    `json:"protocol"`
	Process   string    `json:"process"`
	DetectedAt time.Time `json:"detected_at"`
	Reason    string    `json:"reason"`
}

// PortAnomalyStore detects and stores port anomalies based on time-of-day rules.
type PortAnomalyStore struct {
	mu       sync.RWMutex
	records  []AnomalyRecord
	maxSize  int
	allowedStart int // hour 0-23
	allowedEnd   int // hour 0-23, exclusive
	now      func() time.Time
}

// NewPortAnomalyStore creates a store that flags ports seen outside [allowedStart, allowedEnd).
func NewPortAnomalyStore(allowedStart, allowedEnd, maxSize int) *PortAnomalyStore {
	if maxSize <= 0 {
		maxSize = 200
	}
	return &PortAnomalyStore{
		maxSize:      maxSize,
		allowedStart: allowedStart,
		allowedEnd:   allowedEnd,
		now:          time.Now,
	}
}

// Record checks entries and stores anomalies for those outside allowed hours.
func (s *PortAnomalyStore) Record(entries []scanner.Entry) {
	if len(entries) == 0 {
		return
	}
	now := s.now()
	hour := now.Hour()
	inWindow := hour >= s.allowedStart && hour < s.allowedEnd
	if inWindow {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range entries {
		rec := AnomalyRecord{
			Port:       e.Port,
			Address:    e.Address,
			Protocol:   e.Protocol.String(),
			Process:    e.Process,
			DetectedAt: now,
			Reason:     "port seen outside allowed hours",
		}
		s.records = append(s.records, rec)
		if len(s.records) > s.maxSize {
			s.records = s.records[len(s.records)-s.maxSize:]
		}
	}
}

// Snapshot returns a copy of all recorded anomalies.
func (s *PortAnomalyStore) Snapshot() []AnomalyRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]AnomalyRecord, len(s.records))
	copy(out, s.records)
	return out
}

// NewPortAnomalyAPI returns an HTTP handler serving anomaly records as JSON.
func NewPortAnomalyAPI(store *PortAnomalyStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		records := store.Snapshot()
		if records == nil {
			records = []AnomalyRecord{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(records)
	})
}
