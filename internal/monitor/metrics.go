package monitor

import (
	"sync"
	"time"
)

// Metrics tracks runtime statistics for the monitor daemon.
type Metrics struct {
	mu           sync.RWMutex
	ScanCount     int
	AlertCount    int
	ErrorCount    int
	LastScanTime  time.Time
	LastAlertTime time.Time
	UpSince       time.Time
}

// NewMetrics creates a new Metrics instance with UpSince set to now.
func NewMetrics() *Metrics {
	return &Metrics{
		UpSince: time.Now(),
	}
}

// RecordScan increments the scan counter and updates LastScanTime.
func (m *Metrics) RecordScan() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ScanCount++
	m.LastScanTime = time.Now()
}

// RecordAlert increments the alert counter and updates LastAlertTime.
func (m *Metrics) RecordAlert() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AlertCount++
	m.LastAlertTime = time.Now()
}

// RecordError increments the error counter.
func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorCount++
}

// Snapshot returns a copy of the current metrics.
func (m *Metrics) Snapshot() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return Metrics{
		ScanCount:     m.ScanCount,
		AlertCount:    m.AlertCount,
		ErrorCount:    m.ErrorCount,
		LastScanTime:  m.LastScanTime,
		LastAlertTime: m.LastAlertTime,
		UpSince:       m.UpSince,
	}
}

// UptimeDuration returns how long the monitor has been running.
func (m *Metrics) UptimeDuration() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return time.Since(m.UpSince)
}
