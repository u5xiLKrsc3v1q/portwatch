package monitor

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"

	"github.com/iamcalledrob/portwatch/internal/scanner"
)

// PortRiskEntry holds a risk score and reason for a port.
type PortRiskEntry struct {
	Port     uint16  `json:"port"`
	Protocol string  `json:"protocol"`
	Score    float64 `json:"score"`
	Reason   string  `json:"reason"`
}

// PortRiskStore tracks risk scores per port entry.
type PortRiskStore struct {
	mu      sync.RWMutex
	entries map[string]PortRiskEntry
	classifier *scanner.Classifier
}

func NewPortRiskStore(c *scanner.Classifier) *PortRiskStore {
	return &PortRiskStore{
		entries:    make(map[string]PortRiskEntry),
		classifier: c,
	}
}

func (s *PortRiskStore) Record(e scanner.Entry) {
	if s.classifier == nil {
		return
	}
	sev := s.classifier.Classify(e)
	score, reason := severityToScore(sev)
	key := e.Key()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = PortRiskEntry{
		Port:     e.Port,
		Protocol: e.Protocol.String(),
		Score:    score,
		Reason:   reason,
	}
}

func (s *PortRiskStore) Snapshot() []PortRiskEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]PortRiskEntry, 0, len(s.entries))
	for _, v := range s.entries {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

func severityToScore(sev scanner.Severity) (float64, string) {
	switch sev {
	case scanner.SeverityHigh:
		return 1.0, "well-known privileged port"
	case scanner.SeverityMedium:
		return 0.5, "privileged port"
	default:
		return 0.1, "high/ephemeral port"
	}
}

// NewPortRiskAPI returns an HTTP handler for the risk endpoint.
func NewPortRiskAPI(s *PortRiskStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.Snapshot())
	})
}
