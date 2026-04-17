package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/jwhittle933/portwatch/internal/scanner"
)

// ChangeSummary aggregates added/removed entries over a time window.
type ChangeSummary struct {
	Window    time.Duration
	Added     []scanner.Entry
	Removed   []scanner.Entry
	RecordedAt time.Time
}

// NewChangeSummary builds a ChangeSummary from diff results.
func NewChangeSummary(added, removed []scanner.Entry, window time.Duration) ChangeSummary {
	return ChangeSummary{
		Window:     window,
		Added:      added,
		Removed:    removed,
		RecordedAt: time.Now(),
	}
}

// HasChanges returns true if there are any added or removed entries.
func (s ChangeSummary) HasChanges() bool {
	return len(s.Added) > 0 || len(s.Removed) > 0
}

// String returns a human-readable summary.
func (s ChangeSummary) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "ChangeSummary [%s window, %s]:\n", s.Window, s.RecordedAt.Format(time.RFC3339))
	if len(s.Added) > 0 {
		b.WriteString("  Added:\n")
		for _, e := range s.Added {
			fmt.Fprintf(&b, "    + %s\n", formatEntry(e))
		}
	}
	if len(s.Removed) > 0 {
		b.WriteString("  Removed:\n")
		for _, e := range s.Removed {
			fmt.Fprintf(&b, "    - %s\n", formatEntry(e))
		}
	}
	return b.String()
}

// TotalChanges returns the total count of changes.
func (s ChangeSummary) TotalChanges() int {
	return len(s.Added) + len(s.Removed)
}
