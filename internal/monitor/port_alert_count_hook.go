package monitor

import "github.com/ben/portwatch/internal/scanner"

// PortAlertCountHook records an alert count increment for every added entry
// in an AlertEvent. It is wired into the scan cycle's post-alert hooks.
type PortAlertCountHook struct {
	store *PortAlertCountStore
}

// NewPortAlertCountHook creates a hook backed by the given store.
// If store is nil the hook is a no-op.
func NewPortAlertCountHook(store *PortAlertCountStore) *PortAlertCountHook {
	return &PortAlertCountHook{store: store}
}

// OnAlert is called after an alert is dispatched. It increments the count for
// every newly added port that triggered the alert.
func (h *PortAlertCountHook) OnAlert(event AlertEvent) {
	if h.store == nil {
		return
	}
	for _, e := range event.Added {
		h.store.Record(entryKey(e))
	}
}

// entryKey returns a stable string key for a scanner.Entry.
func entryKey(e scanner.Entry) string {
	return e.Key()
}
