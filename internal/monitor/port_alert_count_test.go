package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPortAlertCountStore_InitiallyEmpty(t *testing.T) {
	s := NewPortAlertCountStore()
	if got := s.Count("tcp:80"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
	if snap := s.Snapshot(); len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(snap))
	}
}

func TestPortAlertCountStore_RecordAndCount(t *testing.T) {
	s := NewPortAlertCountStore()
	s.Record("tcp:80")
	s.Record("tcp:80")
	s.Record("udp:53")
	if got := s.Count("tcp:80"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
	if got := s.Count("udp:53"); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestPortAlertCountStore_Snapshot_SortedDescending(t *testing.T) {
	s := NewPortAlertCountStore()
	s.Record("tcp:443")
	s.Record("tcp:80")
	s.Record("tcp:80")
	s.Record("tcp:80")
	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	if snap[0].Key != "tcp:80" {
		t.Fatalf("expected tcp:80 first, got %s", snap[0].Key)
	}
	if snap[0].Count != 3 {
		t.Fatalf("expected count 3, got %d", snap[0].Count)
	}
}

func TestPortAlertCountStore_Reset(t *testing.T) {
	s := NewPortAlertCountStore()
	s.Record("tcp:80")
	s.Reset()
	if got := s.Count("tcp:80"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestPortAlertCountStore_LastAt_Set(t *testing.T) {
	s := NewPortAlertCountStore()
	before := time.Now()
	s.Record("tcp:22")
	after := time.Now()
	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatal("expected 1 entry")
	}
	if snap[0].LastAt.Before(before) || snap[0].LastAt.After(after) {
		t.Fatalf("LastAt %v out of expected range [%v, %v]", snap[0].LastAt, before, after)
	}
}

func TestPortAlertCountAPI_Get(t *testing.T) {
	s := NewPortAlertCountStore()
	s.Record("tcp:8080")
	h := NewPortAlertCountAPI(s)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []PortAlertCountEntry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) != 1 || entries[0].Key != "tcp:8080" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

func TestPortAlertCountAPI_MethodNotAllowed(t *testing.T) {
	s := NewPortAlertCountStore()
	h := NewPortAlertCountAPI(s)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
