package monitor

import (
	"github.com/user/portwatch/internal/scanner"
)

// PortFirstSeenHook records newly added entries into the PortFirstSeenStore.
type PortFirstSeenHook struct {
	store *PortFirstSeenStore
}

// NewPortFirstSeenHook creates a hook backed by the given store.
func NewPortFirstSeenHook(store *PortFirstSeenStore) *PortFirstSeenHook {
	return &PortFirstSeenHook{store: store}
}

// OnScan is called after each scan cycle with the current list of active entries.
func (h *PortFirstSeenHook) OnScan(entries []scanner.Entry) {
	if h.store == nil {
		return
	}
	for _, e := range entries {
		h.store.Record(e)
	}
}
