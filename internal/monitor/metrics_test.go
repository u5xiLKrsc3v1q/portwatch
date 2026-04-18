package monitor

import (
	"testing"
	"time"
)

func TestNewMetrics_InitialValues(t *testing.T) {
	before := time.Now()
	m := NewMetrics()
	after := time.Now()

	if m.ScanCount != 0 || m.AlertCount != 0 || m.ErrorCount != 0 {
		t.Error("expected zero counters on init")
	}
	if m.UpSince.Before(before) || m.UpSince.After(after) {
		t.Error("UpSince not set correctly")
	}
}

func TestMetrics_RecordScan(t *testing.T) {
	m := NewMetrics()
	m.RecordScan()
	m.RecordScan()
	snap := m.Snapshot()
	if snap.ScanCount != 2 {
		t.Errorf("expected ScanCount=2, got %d", snap.ScanCount)
	}
	if snap.LastScanTime.IsZero() {
		t.Error("expected LastScanTime to be set")
	}
}

func TestMetrics_RecordAlert(t *testing.T) {
	m := NewMetrics()
	m.RecordAlert()
	snap := m.Snapshot()
	if snap.AlertCount != 1 {
		t.Errorf("expected AlertCount=1, got %d", snap.AlertCount)
	}
	if snap.LastAlertTime.IsZero() {
		t.Error("expected LastAlertTime to be set")
	}
}

func TestMetrics_RecordError(t *testing.T) {
	m := NewMetrics()
	m.RecordError()
	m.RecordError()
	m.RecordError()
	if m.Snapshot().ErrorCount != 3 {
		t.Errorf("expected ErrorCount=3, got %d", m.Snapshot().ErrorCount)
	}
}

func TestMetrics_Snapshot_IsCopy(t *testing.T) {
	m := NewMetrics()
	m.RecordScan()
	snap := m.Snapshot()
	m.RecordScan()
	if snap.ScanCount != 1 {
		t.Error("snapshot should not reflect changes after capture")
	}
}

func TestMetrics_UptimeDuration(t *testing.T) {
	m := NewMetrics()
	time.Sleep(5 * time.Millisecond)
	up := m.UptimeDuration()
	if up < 5*time.Millisecond {
		t.Errorf("expected uptime >= 5ms, got %v", up)
	}
}
