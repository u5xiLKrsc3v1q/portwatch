package monitor

import "github.com/iamcalledrob/portwatch/internal/scanner"

// PortRiskHook records risk scores for newly added entries after each scan.
type PortRiskHook struct {
	store *PortRiskStore
}

func NewPortRiskHook(store *PortRiskStore) *PortRiskHook {
	return &PortRiskHook{store: store}
}

func (h *PortRiskHook) OnScan(added, _ []scanner.Entry) {
	for _, e := range added {
		h.store.Record(e)
	}
}
