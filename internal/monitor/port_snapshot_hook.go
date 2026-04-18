package monitor

import "github.com/rgst-io/portwatch/internal/scanner"

// PortSnapshotHook updates the PortSnapshotStore after each scan cycle.
type PortSnapshotHook struct {
	store *PortSnapshotStore
}

func NewPortSnapshotHook(store *PortSnapshotStore) *PortSnapshotHook {
	return &PortSnapshotHook{store: store}
}

func (h *PortSnapshotHook) OnScan(entries []scanner.Entry) {
	h.store.Update(entries)
}
