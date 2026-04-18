package scanner

// Snapshot holds a set of port entries keyed by a stable string key.
type Snapshot map[string]PortEntry

// Key returns a unique string key for a PortEntry.
func (e PortEntry) Key() string {
	return e.Protocol + "|" + e.LocalAddr
}

// NewSnapshot builds a Snapshot from a slice of PortEntry.
func NewSnapshot(entries []PortEntry) Snapshot {
	s := make(Snapshot, len(entries))
	for _, e := range entries {
		s[e.Key()] = e
	}
	return s
}

// DiffResult contains newly appeared and disappeared port entries.
type DiffResult struct {
	Added   []PortEntry
	Removed []PortEntry
}

// Diff compares a previous snapshot against a current one.
func Diff(prev, curr Snapshot) DiffResult {
	var result DiffResult

	for key, entry := range curr {
		if _, exists := prev[key]; !exists {
			result.Added = append(result.Added, entry)
		}
	}

	for key, entry := range prev {
		if _, exists := curr[key]; !exists {
			result.Removed = append(result.Removed, entry)
		}
	}

	return result
}

// HasChanges returns true if the diff contains any additions or removals.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// Equal reports whether two snapshots contain the same set of port entries.
func (s Snapshot) Equal(other Snapshot) bool {
	if len(s) != len(other) {
		return false
	}
	for key := range s {
		if _, exists := other[key]; !exists {
			return false
		}
	}
	return true
}
