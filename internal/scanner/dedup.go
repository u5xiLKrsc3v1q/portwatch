package scanner

import "sync"

// Deduplicator suppresses repeated alerts for the same entry within a
// configurable window. Once an entry has been seen it is held until
// Expire is called (e.g. on each scan cycle) so that the same port
// binding does not fire a notification on every tick.
type Deduplicator struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

// NewDeduplicator returns an initialised Deduplicator.
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{
		seen: make(map[string]struct{}),
	}
}

// IsDuplicate returns true if the entry has already been reported since
// the last call to Reset. It also marks the entry as seen.
func (d *Deduplicator) IsDuplicate(e Entry) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	k := e.Key()
	if _, ok := d.seen[k]; ok {
		return true
	}
	d.seen[k] = struct{}{}
	return false
}

// Filter returns only those entries that have not been seen before,
// marking each returned entry as seen.
func (d *Deduplicator) Filter(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !d.IsDuplicate(e) {
			out = append(out, e)
		}
	}
	return out
}

// Reset clears the seen set so all entries will be treated as new again.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]struct{})
}

// Len returns the number of currently tracked entries.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
