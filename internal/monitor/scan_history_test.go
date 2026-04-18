package monitor

import (
	"errors"
	"testing"
	"time"
)

func makeScanRecord(added, removed, total int, err error) ScanRecord {
	return ScanRecord{
		Timestamp: time.Now(),
		Added:     added,
		Removed:   removed,
		Total:     total,
		Error:     err,
	}
}

func TestScanHistory_DefaultMaxSize(t *testing.T) {
	h := NewScanHistory(0)
	if h.maxSize != 100 {
		t.Errorf("expected default maxSize 100, got %d", h.maxSize)
	}
}

func TestScanHistory_AppendAndLen(t *testing.T) {
	h := NewScanHistory(10)
	h.Record(makeScanRecord(1, 0, 5, nil))
	h.Record(makeScanRecord(0, 1, 4, nil))
	if h.Len() != 2 {
		t.Errorf("expected 2 records, got %d", h.Len())
	}
}

func TestScanHistory_Entries_Snapshot(t *testing.T) {
	h := NewScanHistory(10)
	h.Record(makeScanRecord(3, 0, 3, nil))
	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if entries[0].Added != 3 {
		t.Errorf("expected Added=3")
	}
}

func TestScanHistory_EvictsOldest(t *testing.T) {
	h := NewScanHistory(3)
	for i := 0; i < 5; i++ {
		h.Record(makeScanRecord(i, 0, i, nil))
	}
	if h.Len() != 3 {
		t.Errorf("expected 3 after eviction, got %d", h.Len())
	}
	entries := h.Entries()
	if entries[0].Added != 2 {
		t.Errorf("expected oldest remaining Added=2, got %d", entries[0].Added)
	}
}

func TestScanHistory_LastError_None(t *testing.T) {
	h := NewScanHistory(10)
	if h.LastError() != nil {
		t.Error("expected nil error on empty history")
	}
	h.Record(makeScanRecord(0, 0, 0, nil))
	if h.LastError() != nil {
		t.Error("expected nil error")
	}
}

func TestScanHistory_LastError_Set(t *testing.T) {
	h := NewScanHistory(10)
	sentinel := errors.New("scan failed")
	h.Record(makeScanRecord(0, 0, 0, sentinel))
	if h.LastError() != sentinel {
		t.Errorf("expected sentinel error, got %v", h.LastError())
	}
}
