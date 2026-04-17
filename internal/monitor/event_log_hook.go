package monitor

import (
	"github.com/sgreben/portwatch/internal/scanner"
)

// EventLogHook wires an EventLog into the scan cycle as a side-effect observer.
// It records every alert event without interfering with the notifier pipeline.
type EventLogHook struct {
	log *EventLog
}

// NewEventLogHook creates a hook backed by the given EventLog.
func NewEventLogHook(log *EventLog) *EventLogHook {
	return &EventLogHook{log: log}
}

// Record appends the added/removed entries to the event log.
func (h *EventLogHook) Record(added, removed []scanner.Entry) {
	if h.log == nil {
		return
	}
	if len(added) == 0 && len(removed) == 0 {
		return
	}
	h.log.Append(added, removed)
}
