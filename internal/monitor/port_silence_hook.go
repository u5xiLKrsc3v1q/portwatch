package monitor

import (
	"log"

	"github.com/user/portwatch/internal/scanner"
)

// PortSilenceHook filters out silenced ports from added entries before alerting.
type PortSilenceHook struct {
	store *PortSilenceStore
}

// NewPortSilenceHook creates a PortSilenceHook backed by the given store.
func NewPortSilenceHook(store *PortSilenceStore) *PortSilenceHook {
	return &PortSilenceHook{store: store}
}

// FilterAdded removes any added entries that are currently silenced.
func (h *PortSilenceHook) FilterAdded(entries []scanner.Entry) []scanner.Entry {
	if h.store == nil {
		return entries
	}
	out := entries[:0]
	for _, e := range entries {
		if h.store.IsSilenced(e) {
			log.Printf("[portsilence] suppressing silenced port %s", e.Key())
			continue
		}
		out = append(out, e)
	}
	return out
}
