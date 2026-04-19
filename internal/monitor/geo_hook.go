package monitor

import (
	"log"

	"github.com/robherley/portwatch/internal/scanner"
)

// GeoHook integrates GeoFilter into the scan cycle pipeline.
// It logs and optionally strips entries from denied countries.
type GeoHook struct {
	filter *GeoFilter
	strip  bool
}

// NewGeoHook creates a hook that uses the provided GeoFilter.
// When strip is true, denied entries are removed from Added lists.
func NewGeoHook(f *GeoFilter, strip bool) *GeoHook {
	return &GeoHook{filter: f, strip: strip}
}

// OnScan inspects added entries and optionally removes geo-denied ones.
func (h *GeoHook) OnScan(added, removed []scanner.Entry) ([]scanner.Entry, []scanner.Entry) {
	if h == nil || h.filter == nil {
		return added, removed
	}
	filtered := added[:0:len(added)]
	for _, e := range added {
		if h.filter.IsDenied(e.Address) {
			log.Printf("[geo] denied entry addr=%s port=%d", e.Address, e.Port)
			if h.strip {
				continue
			}
		}
		filtered = append(filtered, e)
	}
	return filtered, removed
}
