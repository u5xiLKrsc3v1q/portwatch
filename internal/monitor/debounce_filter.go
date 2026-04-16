package monitor

import (
	"time"

	"github.com/danvolchek/portwatch/internal/scanner"
)

// DebounceFilter wraps a Debouncer to suppress repeated added-entry alerts.
type DebounceFilter struct {
	debouncer *Debouncer
}

// NewDebounceFilter creates a DebounceFilter with the given quiet window.
func NewDebounceFilter(window time.Duration) *DebounceFilter {
	return &DebounceFilter{debouncer: NewDebouncer(window)}
}

// Filter removes added entries that are within the debounce window.
// Removed entries always pass through.
func (f *DebounceFilter) Filter(added, removed []scanner.Entry) ([]scanner.Entry, []scanner.Entry) {
	if f == nil || f.debouncer == nil {
		return added, removed
	}
	var allowed []scanner.Entry
	for _, e := range added {
		if f.debouncer.Allow(e.Key()) {
			allowed = append(allowed, e)
		}
	}
	return allowed, removed
}
