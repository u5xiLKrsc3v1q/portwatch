package monitor

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// PortNote holds a user-defined annotation for a port key.
type PortNote struct {
	Note      string    `json:"note"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PortNoteStore stores user annotations keyed by port identifier.
type PortNoteStore struct {
	mu    sync.RWMutex
	notes map[string]PortNote
}

// NewPortNoteStore creates an empty PortNoteStore.
func NewPortNoteStore() *PortNoteStore {
	return &PortNoteStore{
		notes: make(map[string]PortNote),
	}
}

// Set stores or updates a note for the given key.
func (s *PortNoteStore) Set(key, note string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[key] = PortNote{
		Note:      note,
		UpdatedAt: time.Now(),
	}
}

// Get retrieves the note for the given key. Returns false if not found.
func (s *PortNoteStore) Get(key string) (PortNote, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.notes[key]
	return n, ok
}

// Delete removes the note for the given key.
func (s *PortNoteStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.notes, key)
}

// Snapshot returns a copy of all notes.
func (s *PortNoteStore) Snapshot() map[string]PortNote {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]PortNote, len(s.notes))
	for k, v := range s.notes {
		out[k] = v
	}
	return out
}

// NewPortNoteAPI returns an http.Handler that exposes GET/POST/DELETE for port notes.
func NewPortNoteAPI(store *PortNoteStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(store.Snapshot())
		case http.MethodPost:
			var req struct {
				Key  string `json:"key"`
				Note string `json:"note"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Key == "" {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			store.Set(req.Key, req.Note)
			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			key := r.URL.Query().Get("key")
			if key == "" {
				http.Error(w, "missing key", http.StatusBadRequest)
				return
			}
			store.Delete(key)
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}
