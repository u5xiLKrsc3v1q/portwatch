package monitor

import (
	"sync"
	"time"
)

// ScanRecord holds the result summary of a single scan cycle.
type ScanRecord struct {
	Timestamp time.Time
	Added     int
	Removed   int
	Total     int
	Error     error
}

// ScanHistory keeps a bounded ring of recent scan records.
type ScanHistory struct {
	mu      sync.Mutex
	records []ScanRecord
	maxSize int
}

// NewScanHistory creates a ScanHistory with the given capacity.
func NewScanHistory(maxSize int) *ScanHistory {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &ScanHistory{maxSize: maxSize}
}

// Record appends a new scan result, evicting the oldest if at capacity.
func (h *ScanHistory) Record(r ScanRecord) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.records) >= h.maxSize {
		h.records = h.records[1:]
	}
	h.records = append(h.records, r)
}

// Entries returns a snapshot of all stored records.
func (h *ScanHistory) Entries() []ScanRecord {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]ScanRecord, len(h.records))
	copy(out, h.records)
	return out
}

// Len returns the current number of stored records.
func (h *ScanHistory) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.records)
}

// LastError returns the error from the most recent scan, or nil.
func (h *ScanHistory) LastError() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.records) == 0 {
		return nil
	}
	return h.records[len(h.records)-1].Error
}
