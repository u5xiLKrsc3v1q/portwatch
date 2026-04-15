package scanner

import (
	"sync"
	"time"
)

// StateStore tracks the last known snapshot and when it was captured.
type StateStore struct {
	mu       sync.RWMutex
	snapshot *Snapshot
	updatedAt time.Time
}

// NewStateStore creates an empty StateStore.
func NewStateStore() *StateStore {
	return &StateStore{}
}

// Set atomically replaces the current snapshot.
func (s *StateStore) Set(snap *Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshot = snap
	s.updatedAt = time.Now()
}

// Get returns the current snapshot and the time it was last set.
// Returns nil, zero-time if no snapshot has been stored yet.
func (s *StateStore) Get() (*Snapshot, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshot, s.updatedAt
}

// HasSnapshot reports whether a snapshot has been stored.
func (s *StateStore) HasSnapshot() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshot != nil
}

// UpdateAndDiff stores the new snapshot and returns the diff against the
// previous one. If no previous snapshot exists the diff will show all
// entries in newSnap as Added.
func (s *StateStore) UpdateAndDiff(newSnap *Snapshot) Changes {
	s.mu.Lock()
	defer s.mu.Unlock()

	prev := s.snapshot
	s.snapshot = newSnap
	s.updatedAt = time.Now()

	if prev == nil {
		// First run — treat everything as newly added.
		var changes Changes
		for _, e := range newSnap.Entries {
			changes.Added = append(changes.Added, e)
		}
		return changes
	}

	return Diff(prev, newSnap)
}
