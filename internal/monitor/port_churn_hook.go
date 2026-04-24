package monitor

import "github.com/user/portwatch/internal/scanner"

// PortChurnHook records churn events into a PortChurnStore whenever ports are
// added or removed during a scan cycle.
type PortChurnHook struct {
	store *PortChurnStore
}

// NewPortChurnHook returns a hook that feeds diff events into store.
func NewPortChurnHook(store *PortChurnStore) *PortChurnHook {
	return &PortChurnHook{store: store}
}

// OnScan is called after each scan cycle with added and removed entries.
func (h *PortChurnHook) OnScan(added, removed []scanner.Entry) {
	if h.store == nil {
		return
	}
	for _, e := range added {
		h.store.Record(e)
	}
	for _, e := range removed {
		h.store.Record(e)
	}
}
