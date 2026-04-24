package monitor

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// PortOwnerStore maps ports to their associated process names.
type PortOwnerStore struct {
	mu    sync.RWMutex
	owners map[uint16]string
}

// NewPortOwnerStore creates an empty PortOwnerStore.
func NewPortOwnerStore() *PortOwnerStore {
	return &PortOwnerStore{
		owners: make(map[uint16]string),
	}
}

// Record stores the process name for a given port. If the process name is
// empty the entry is removed.
func (s *PortOwnerStore) Record(port uint16, process string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if process == "" {
		delete(s.owners, port)
		return
	}
	s.owners[port] = process
}

// Get returns the process name for a port and whether it was found.
func (s *PortOwnerStore) Get(port uint16) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.owners[port]
	return v, ok
}

// Snapshot returns a copy of the current owner map.
func (s *PortOwnerStore) Snapshot() map[uint16]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[uint16]string, len(s.owners))
	for k, v := range s.owners {
		out[k] = v
	}
	return out
}

// NewPortOwnerAPI returns an http.Handler that serves the owner snapshot.
func NewPortOwnerAPI(store *PortOwnerStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Snapshot())
	})
}

// ownerEntryKey returns the port number from an entry as a map key.
func ownerEntryKey(e scanner.Entry) uint16 {
	return e.Port
}
