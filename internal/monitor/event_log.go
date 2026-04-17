package monitor

import (
	"sync"
	"time"

	"github.com/sgreben/portwatch/internal/scanner"
)

// EventLogEntry records a single alert event with a timestamp.
type EventLogEntry struct {
	Time    time.Time
	Added   []scanner.Entry
	Removed []scanner.Entry
}

// EventLog stores a bounded in-memory history of alert events.
type EventLog struct {
	mu      sync.RWMutex
	entries []EventLogEntry
	maxSize int
}

// NewEventLog creates an EventLog with the given maximum capacity.
func NewEventLog(maxSize int) *EventLog {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &EventLog{maxSize: maxSize}
}

// Append adds a new entry to the log, evicting the oldest if at capacity.
func (l *EventLog) Append(added, removed []scanner.Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry := EventLogEntry{
		Time:    time.Now(),
		Added:   added,
		Removed: removed,
	}
	if len(l.entries) >= l.maxSize {
		l.entries = l.entries[1:]
	}
	l.entries = append(l.entries, entry)
}

// Entries returns a snapshot of all log entries.
func (l *EventLog) Entries() []EventLogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]EventLogEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Len returns the current number of entries.
func (l *EventLog) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries)
}
