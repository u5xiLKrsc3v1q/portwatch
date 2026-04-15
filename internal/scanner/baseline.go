package scanner

import (
	"encoding/json"
	"os"
	"sync"
)

// Baseline persists a known-good set of port entries to disk so that
// portwatch can distinguish "already present at startup" listeners from
// truly new ones across restarts.
type Baseline struct {
	mu      sync.RWMutex
	entries map[string]Entry
	path    string
}

// NewBaseline loads an existing baseline file from path, or returns an
// empty baseline if the file does not yet exist.
func NewBaseline(path string) (*Baseline, error) {
	b := &Baseline{
		entries: make(map[string]Entry),
		path:    path,
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return b, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, e := range list {
		b.entries[e.Key()] = e
	}
	return b, nil
}

// Contains reports whether the entry is part of the baseline.
func (b *Baseline) Contains(e Entry) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.entries[e.Key()]
	return ok
}

// Save writes the current snapshot of entries as the new baseline.
func (b *Baseline) Save(entries []Entry) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = make(map[string]Entry, len(entries))
	for _, e := range entries {
		b.entries[e.Key()] = e
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

// Entries returns a copy of all baseline entries.
func (b *Baseline) Entries() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		out = append(out, e)
	}
	return out
}
