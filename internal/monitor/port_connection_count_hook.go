package monitor

import "github.com/user/portwatch/internal/scanner"

// PortConnectionCountHook records active entries into a PortConnectionCountStore
// on every scan cycle.
type PortConnectionCountHook struct {
	store *PortConnectionCountStore
}

// NewPortConnectionCountHook creates a hook that feeds scan results into the store.
func NewPortConnectionCountHook(store *PortConnectionCountStore) *PortConnectionCountHook {
	return &PortConnectionCountHook{store: store}
}

// OnScan is called after each scan with the current set of active entries.
func (h *PortConnectionCountHook) OnScan(entries []scanner.Entry) {
	if h.store == nil {
		return
	}
	h.store.Record(entries)
}
