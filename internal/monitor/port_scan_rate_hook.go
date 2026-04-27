package monitor

import (
	"time"

	"github.com/wiring/portwatch/internal/scanner"
)

// PortScanRateHook records port scan frequency into a PortScanRateStore
// after each scan cycle completes.
type PortScanRateHook struct {
	store *PortScanRateStore
}

// NewPortScanRateHook creates a hook that feeds scan results into the store.
func NewPortScanRateHook(store *PortScanRateStore) *PortScanRateHook {
	return &PortScanRateHook{store: store}
}

// OnScan is called after every scan with the full current set of entries.
func (h *PortScanRateHook) OnScan(entries []scanner.Entry) {
	if h.store == nil {
		return
	}
	now := time.Now()
	for _, e := range entries {
		h.store.Record(e.Port, e.Protocol.String(), now)
	}
}
