package monitor

import (
	"time"

	"github.com/rgst-io/portwatch/internal/scanner"
)

// PortLifetimeHook updates the PortLifetimeStore on every scan cycle.
type PortLifetimeHook struct {
	store *PortLifetimeStore
	now   func() time.Time
}

// NewPortLifetimeHook creates a hook that feeds scan results into the store.
func NewPortLifetimeHook(store *PortLifetimeStore) *PortLifetimeHook {
	return &PortLifetimeHook{store: store, now: time.Now}
}

// OnScan is called after each scan with the current set of active entries.
func (h *PortLifetimeHook) OnScan(entries []scanner.Entry) {
	now := h.now()
	seen := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		key := e.Key()
		seen[key] = struct{}{}
		h.store.Record(key, now)
	}
	// Remove entries that are no longer present.
	for _, r := range h.store.Snapshot() {
		if _, ok := seen[r.Key]; !ok {
			h.store.Remove(r.Key)
		}
	}
}
