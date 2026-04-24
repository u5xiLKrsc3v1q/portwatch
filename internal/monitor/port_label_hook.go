package monitor

import "github.com/jwhittle933/portwatch/internal/scanner"

// PortLabelHook annotates scan entries with their label from the store.
// It does not modify the entry itself but can be used to enrich alert events.
type PortLabelHook struct {
	store *PortLabelStore
}

// NewPortLabelHook creates a PortLabelHook backed by the given store.
func NewPortLabelHook(store *PortLabelStore) *PortLabelHook {
	return &PortLabelHook{store: store}
}

// LabelFor returns the human-readable label for the given entry's port,
// or an empty string if none is registered.
func (h *PortLabelHook) LabelFor(e scanner.Entry) string {
	if h == nil || h.store == nil {
		return ""
	}
	return h.store.Get(e.Port)
}

// AnnotateAdded returns a map of port -> label for all added entries
// that have a known label.
func (h *PortLabelHook) AnnotateAdded(entries []scanner.Entry) map[uint16]string {
	out := make(map[uint16]string)
	if h == nil || h.store == nil {
		return out
	}
	for _, e := range entries {
		if label := h.store.Get(e.Port); label != "" {
			out[e.Port] = label
		}
	}
	return out
}
