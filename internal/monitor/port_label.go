package monitor

import (
	"encoding/json"
	"net/http"
	"sync"
)

// PortLabelStore maps port numbers to human-readable labels.
type PortLabelStore struct {
	mu     sync.RWMutex
	labels map[uint16]string
}

// NewPortLabelStore creates a PortLabelStore pre-populated with well-known port labels.
func NewPortLabelStore() *PortLabelStore {
	s := &PortLabelStore{
		labels: map[uint16]string{
			22:   "SSH",
			25:   "SMTP",
			53:   "DNS",
			80:   "HTTP",
			110:  "POP3",
			143:  "IMAP",
			443:  "HTTPS",
			3306: "MySQL",
			5432: "PostgreSQL",
			6379: "Redis",
			8080: "HTTP-Alt",
			27017: "MongoDB",
		},
	}
	return s
}

// Set adds or updates a label for the given port.
func (s *PortLabelStore) Set(port uint16, label string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.labels[port] = label
}

// Get returns the label for the given port, or an empty string if not found.
func (s *PortLabelStore) Get(port uint16) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.labels[port]
}

// Snapshot returns a copy of all labels.
func (s *PortLabelStore) Snapshot() map[uint16]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[uint16]string, len(s.labels))
	for k, v := range s.labels {
		out[k] = v
	}
	return out
}

// NewPortLabelAPI returns an HTTP handler that exposes port labels.
func NewPortLabelAPI(store *PortLabelStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(store.Snapshot())
	})
}
