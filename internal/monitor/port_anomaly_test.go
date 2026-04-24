package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeAnomalyEntry(port uint16, proto scanner.Protocol) scanner.Entry {
	return scanner.Entry{
		Port:     port,
		Address:  "0.0.0.0",
		Protocol: proto,
		Process:  "testd",
	}
}

func fixedTime(hour int) func() time.Time {
	return func() time.Time {
		return time.Date(2024, 1, 15, hour, 0, 0, 0, time.UTC)
	}
}

func TestPortAnomalyStore_NoAnomalyInWindow(t *testing.T) {
	s := NewPortAnomalyStore(8, 20, 100)
	s.now = fixedTime(10) // inside window
	s.Record([]scanner.Entry{makeAnomalyEntry(8080, scanner.TCP)})
	if got := s.Snapshot(); len(got) != 0 {
		t.Fatalf("expected 0 anomalies inside window, got %d", len(got))
	}
}

func TestPortAnomalyStore_AnomalyOutsideWindow(t *testing.T) {
	s := NewPortAnomalyStore(8, 20, 100)
	s.now = fixedTime(3) // outside window
	s.Record([]scanner.Entry{makeAnomalyEntry(22, scanner.TCP)})
	got := s.Snapshot()
	if len(got) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(got))
	}
	if got[0].Port != 22 {
		t.Errorf("expected port 22, got %d", got[0].Port)
	}
	if got[0].Reason == "" {
		t.Error("expected non-empty reason")
	}
}

func TestPortAnomalyStore_EvictsOldest(t *testing.T) {
	s := NewPortAnomalyStore(8, 20, 3)
	s.now = fixedTime(2)
	for i := uint16(1); i <= 5; i++ {
		s.Record([]scanner.Entry{makeAnomalyEntry(i, scanner.TCP)})
	}
	got := s.Snapshot()
	if len(got) != 3 {
		t.Fatalf("expected 3 records after eviction, got %d", len(got))
	}
	if got[0].Port != 3 {
		t.Errorf("expected oldest evicted, first port = 3, got %d", got[0].Port)
	}
}

func TestPortAnomalyStore_EmptyRecord(t *testing.T) {
	s := NewPortAnomalyStore(8, 20, 100)
	s.now = fixedTime(2)
	s.Record(nil)
	if len(s.Snapshot()) != 0 {
		t.Error("expected no anomalies for nil input")
	}
}

func TestPortAnomalyAPI_Get(t *testing.T) {
	s := NewPortAnomalyStore(8, 20, 100)
	s.now = fixedTime(1)
	s.Record([]scanner.Entry{makeAnomalyEntry(443, scanner.TCP)})

	handler := NewPortAnomalyAPI(s)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/anomalies", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var records []AnomalyRecord
	if err := json.NewDecoder(rec.Body).Decode(&records); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(records) != 1 || records[0].Port != 443 {
		t.Errorf("unexpected records: %+v", records)
	}
}

func TestPortAnomalyAPI_MethodNotAllowed(t *testing.T) {
	s := NewPortAnomalyStore(8, 20, 100)
	handler := NewPortAnomalyAPI(s)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/anomalies", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
