package monitor

import (
	"net/http"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PortSilenceEntry records a port that has been silenced (suppressed from alerts).
type PortSilenceEntry struct {
	Key       string    `json:"key"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	SilencedAt time.Time `json:"silenced_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// PortSilenceStore tracks manually silenced ports.
type PortSilenceStore struct {
	mu      sync.RWMutex
	entries map[string]PortSilenceEntry
}

// NewPortSilenceStore creates an empty PortSilenceStore.
func NewPortSilenceStore() *PortSilenceStore {
	return &PortSilenceStore{
		entries: make(map[string]PortSilenceEntry),
	}
}

// Silence adds a port to the silence list with an optional expiry duration.
// If ttl is 0, the silence never expires.
func (s *PortSilenceStore) Silence(e scanner.Entry, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := e.Key()
	entry := PortSilenceEntry{
		Key:        key,
		Port:       e.Port,
		Protocol:   e.Protocol.String(),
		SilencedAt: time.Now(),
	}
	if ttl > 0 {
		t := time.Now().Add(ttl)
		entry.ExpiresAt = &t
	}
	s.entries[key] = entry
}

// IsSilenced returns true if the entry is currently silenced.
func (s *PortSilenceStore) IsSilenced(e scanner.Entry) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.entries[e.Key()]
	if !ok {
		return false
	}
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return false
	}
	return true
}

// Unsilence removes a port from the silence list.
func (s *PortSilenceStore) Unsilence(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Snapshot returns a copy of all current silence entries (excluding expired).
func (s *PortSilenceStore) Snapshot() []PortSilenceEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	out := make([]PortSilenceEntry, 0, len(s.entries))
	for _, e := range s.entries {
		if e.ExpiresAt != nil && now.After(*e.ExpiresAt) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// Purge removes all expired silence entries from the store.
// It returns the number of entries removed.
func (s *PortSilenceStore) Purge() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	count := 0
	for key, e := range s.entries {
		if e.ExpiresAt != nil && now.After(*e.ExpiresAt) {
			delete(s.entries, key)
			count++
		}
	}
	return count
}

// NewPortSilenceAPI returns an http.Handler for the silence store.
func NewPortSilenceAPI(store *PortSilenceStore) http.Handler {
	return &portSilenceAPI{store: store}
}

type portSilenceAPI struct{ store *PortSilenceStore }

func (a *portSilenceAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, a.store.Snapshot())
}
