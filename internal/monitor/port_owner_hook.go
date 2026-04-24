package monitor

import "github.com/user/portwatch/internal/scanner"

// PortOwnerHook updates a PortOwnerStore after each scan cycle.
type PortOwnerHook struct {
	store *PortOwnerStore
}

// NewPortOwnerHook creates a PortOwnerHook backed by the given store.
func NewPortOwnerHook(store *PortOwnerStore) *PortOwnerHook {
	return &PortOwnerHook{store: store}
}

// OnScan is called after every scan with the full current entry list.
func (h *PortOwnerHook) OnScan(entries []scanner.Entry) {
	if h.store == nil {
		return
	}
	// Build a fresh set of ports seen this cycle.
	seen := make(map[uint16]struct{}, len(entries))
	for _, e := range entries {
		h.store.Record(e.Port, e.Process)
		seen[e.Port] = struct{}{}
	}
	// Remove ports that are no longer present.
	for port := range h.store.Snapshot() {
		if _, ok := seen[port]; !ok {
			h.store.Record(port, "")
		}
	}
}
