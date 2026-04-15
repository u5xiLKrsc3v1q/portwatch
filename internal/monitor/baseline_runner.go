package monitor

import (
	"log"

	"github.com/iamcalledned/portwatch/internal/scanner"
)

// BaselineManager wraps a Baseline and provides helpers used by the Monitor
// to decide whether a newly detected entry should be suppressed (it was
// present when the baseline was last saved) and to optionally persist a new
// baseline after the first scan.
type BaselineManager struct {
	baseline     *scanner.Baseline
	saveOnce     bool
	savedAlready bool
}

// NewBaselineManager creates a BaselineManager. If saveOnce is true the
// manager will overwrite the baseline file the first time SaveIfNeeded is
// called.
func NewBaselineManager(b *scanner.Baseline, saveOnce bool) *BaselineManager {
	return &BaselineManager{baseline: b, saveOnce: saveOnce}
}

// IsBaseline returns true when the entry is part of the stored baseline and
// should therefore not trigger an alert.
func (bm *BaselineManager) IsBaseline(e scanner.Entry) bool {
	if bm.baseline == nil {
		return false
	}
	return bm.baseline.Contains(e)
}

// FilterAdded removes entries that are already in the baseline from the
// added list, returning only genuinely new listeners.
func (bm *BaselineManager) FilterAdded(added []scanner.Entry) []scanner.Entry {
	if bm.baseline == nil {
		return added
	}
	out := added[:0]
	for _, e := range added {
		if !bm.baseline.Contains(e) {
			out = append(out, e)
		}
	}
	return out
}

// SaveIfNeeded persists entries as the new baseline if saveOnce is enabled
// and the baseline has not yet been saved during this run.
func (bm *BaselineManager) SaveIfNeeded(entries []scanner.Entry) {
	if bm.baseline == nil || !bm.saveOnce || bm.savedAlready {
		return
	}
	if err := bm.baseline.Save(entries); err != nil {
		log.Printf("portwatch: failed to save baseline: %v", err)
		return
	}
	bm.savedAlready = true
	log.Printf("portwatch: baseline saved (%d entries)", len(entries))
}
