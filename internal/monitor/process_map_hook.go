package monitor

import "github.com/user/portwatch/internal/scanner"

// ProcessMapHook updates the ProcessMapStore after each scan cycle.
type ProcessMapHook struct {
	store *ProcessMapStore
}

func NewProcessMapHook(store *ProcessMapStore) *ProcessMapHook {
	return &ProcessMapHook{store: store}
}

func (h *ProcessMapHook) OnScan(entries []scanner.Entry) {
	h.store.Update(entries)
}
