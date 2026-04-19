package monitor

import (
	"sync"
	"time"
)

// AlertRecord captures a single alert event for history tracking.
type AlertRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
	Summary   string    `json:"summary"`
}

// AlertHistory stores a bounded list of past alert events.
type AlertHistory struct {
	mu      sync.Mutex
	records []AlertRecord
	maxSize int
}

// NewAlertHistory creates an AlertHistory with the given capacity.
func NewAlertHistory(maxSize int) *AlertHistory {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &AlertHistory{maxSize: maxSize}
}

// Append adds a new alert record, evicting the oldest if at capacity.
func (h *AlertHistory) Append(r AlertRecord) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.records) >= h.maxSize {
		h.records = h.records[1:]
	}
	h.records = append(h.records, r)
}

// Entries returns a snapshot of all alert records.
func (h *AlertHistory) Entries() []AlertRecord {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]AlertRecord, len(h.records))
	copy(out, h.records)
	return out
}

// Len returns the current number of stored records.
func (h *AlertHistory) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.records)
}
