package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// AlertEvent represents a single port change alert.
type AlertEvent struct {
	Timestamp time.Time
	Added     []scanner.Entry
	Removed   []scanner.Entry
}

// HasChanges returns true if the event contains any added or removed entries.
func (a AlertEvent) HasChanges() bool {
	return len(a.Added) > 0 || len(a.Removed) > 0
}

// Summary returns a human-readable summary of the alert event.
func (a AlertEvent) Summary() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] Port changes detected:\n", a.Timestamp.Format(time.RFC3339)))
	for _, e := range a.Added {
		sb.WriteString(fmt.Sprintf("  + ADDED   %s\n", formatEntry(e)))
	}
	for _, e := range a.Removed {
		sb.WriteString(fmt.Sprintf("  - REMOVED %s\n", formatEntry(e)))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Title returns a short title suitable for a desktop notification.
func (a AlertEvent) Title() string {
	parts := make([]string, 0, 2)
	if len(a.Added) > 0 {
		parts = append(parts, fmt.Sprintf("%d new", len(a.Added)))
	}
	if len(a.Removed) > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", len(a.Removed)))
	}
	return fmt.Sprintf("portwatch: %s listener(s)", strings.Join(parts, ", "))
}
