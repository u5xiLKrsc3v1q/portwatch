package monitor

import (
	"time"

	"github.com/iamcathal/portwatch/internal/scanner"
)

// ScanHistoryHook records each scan cycle result into a ScanHistory.
type ScanHistoryHook struct {
	history *ScanHistory
}

// NewScanHistoryHook creates a hook that writes into the given ScanHistory.
func NewScanHistoryHook(h *ScanHistory) *ScanHistoryHook {
	return &ScanHistoryHook{history: h}
}

// Record builds a ScanRecord from a diff result and appends it to history.
func (s *ScanHistoryHook) Record(added, removed []scanner.Entry, total int, err error) {
	s.history.Record(ScanRecord{
		Timestamp: time.Now(),
		Added:     len(added),
		Removed:   len(removed),
		Total:     total,
		Error:     err,
	})
}
