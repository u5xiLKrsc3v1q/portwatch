package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPortScanRateStore_InitiallyEmpty(t *testing.T) {
	s := NewPortScanRateStore()
	if got := s.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(got))
	}
}

func TestPortScanRateStore_RecordAndCount(t *testing.T) {
	s := NewPortScanRateStore()
	now := time.Now()
	s.Record(80, "tcp", now)
	s.Record(80, "tcp", now.Add(time.Minute))

	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	if snap[0].ScanCount != 2 {
		t.Errorf("expected scan count 2, got %d", snap[0].ScanCount)
	}
	if snap[0].Port != 80 {
		t.Errorf("expected port 80, got %d", snap[0].Port)
	}
}

func TestPortScanRateStore_RatePerMin(t *testing.T) {
	s := NewPortScanRateStore()
	base := time.Now()
	s.Record(443, "tcp", base)
	s.Record(443, "tcp", base.Add(2*time.Minute))

	snap := s.Snapshot()
	if len(snap) == 0 {
		t.Fatal("expected entry")
	}
	if snap[0].RatePerMin <= 0 {
		t.Errorf("expected positive rate, got %f", snap[0].RatePerMin)
	}
}

func TestPortScanRateStore_SortedDescending(t *testing.T) {
	s := NewPortScanRateStore()
	now := time.Now()
	s.Record(22, "tcp", now)
	s.Record(80, "tcp", now)
	s.Record(80, "tcp", now.Add(time.Second))
	s.Record(80, "tcp", now.Add(2*time.Second))

	snap := s.Snapshot()
	if len(snap) < 2 {
		t.Fatal("expected at least 2 entries")
	}
	if snap[0].Port != 80 {
		t.Errorf("expected port 80 first (highest count), got %d", snap[0].Port)
	}
}

func TestPortScanRateAPI_Get(t *testing.T) {
	s := NewPortScanRateStore()
	s.Record(8080, "tcp", time.Now())

	api := NewPortScanRateAPI(s)
	req := httptest.NewRequest(http.MethodGet, "/scan-rate", nil)
	rw := httptest.NewRecorder()
	api.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var result []PortScanRateEntry
	if err := json.NewDecoder(rw.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(result) != 1 || result[0].Port != 8080 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestPortScanRateAPI_MethodNotAllowed(t *testing.T) {
	api := NewPortScanRateAPI(NewPortScanRateStore())
	req := httptest.NewRequest(http.MethodPost, "/scan-rate", nil)
	rw := httptest.NewRecorder()
	api.ServeHTTP(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rw.Code)
	}
}
