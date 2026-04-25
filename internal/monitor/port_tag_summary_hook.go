package monitor

import "github.com/user/portwatch/internal/scanner"

// PortTagSummaryHook records tag counts from each scan cycle.
type PortTagSummaryHook struct {
	store *PortTagSummaryStore
}

// NewPortTagSummaryHook creates a hook that feeds scan results into the store.
func NewPortTagSummaryHook(store *PortTagSummaryStore) *PortTagSummaryHook {
	return &PortTagSummaryHook{store: store}
}

// OnScan is called after each scan with the current set of active entries.
func (h *PortTagSummaryHook) OnScan(entries []scanner.Entry) {
	if h.store == nil {
		return
	}
	h.store.Record(entries)
}
