package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rgst-io/portwatch/internal/scanner"
)

func TestPortLifetimeStore_Record_FirstTime(t *testing.T) {
	s := NewPortLifetimeStore()
	now := time.Now()
	s.Record("tcp:127.0.0.1:8080", now)
	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 record, got %d", len(snap))
	}
	if snap[0].Key != "tcp:127.0.0.1:8080" {
		t.Errorf("unexpected key: %s", snap[0].Key)
	}
	if snap[0].Duration != "0s" {
		t.Errorf("expected 0s duration on first record, got %s", snap[0].Duration)
	}
}

func TestPortLifetimeStore_Record_Idempotent(t *testing.T) {
	s := NewPortLifetimeStore()
	t0 := time.Now()
	t1 := t0.Add(5 * time.Second)
	s.Record("tcp:0.0.0.0:443", t0)
	s.Record("tcp:0.0.0.0:443", t1)
	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 record, got %d", len(snap))
	}
	if snap[0].Duration != "5s" {
		t.Errorf("expected 5s duration, got %s", snap[0].Duration)
	}
	if !snap[0].FirstSeen.Equal(t0) {
		t.Errorf("first_seen should not change on update")
	}
}

func TestPortLifetimeStore_Remove(t *testing.T) {
	s := NewPortLifetimeStore()
	s.Record("tcp:0.0.0.0:22", time.Now())
	s.Remove("tcp:0.0.0.0:22")
	if len(s.Snapshot()) != 0 {
		t.Error("expected empty snapshot after remove")
	}
}

func TestPortLifetimeAPI_Get(t *testing.T) {
	s := NewPortLifetimeStore()
	s.Record("tcp:0.0.0.0:80", time.Now())
	api := NewPortLifetimeAPI(s)
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/lifetime", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var records []PortLifetimeRecord
	if err := json.NewDecoder(rec.Body).Decode(&records); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("expected 1 record in response, got %d", len(records))
	}
}

func TestPortLifetimeAPI_MethodNotAllowed(t *testing.T) {
	api := NewPortLifetimeAPI(NewPortLifetimeStore())
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/lifetime", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestPortLifetimeHook_OnScan_RecordsAndRemoves(t *testing.T) {
	store := NewPortLifetimeStore()
	hook := NewPortLifetimeHook(store)
	fixed := time.Now()
	hook.now = func() time.Time { return fixed }

	e1 := scanner.Entry{Address: "0.0.0.0", Port: 80}
	e2 := scanner.Entry{Address: "0.0.0.0", Port: 443}
	hook.OnScan([]scanner.Entry{e1, e2})
	if len(store.Snapshot()) != 2 {
		t.Fatalf("expected 2 records after first scan")
	}
	// Second scan: e2 disappears.
	hook.OnScan([]scanner.Entry{e1})
	snap := store.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 record after removal, got %d", len(snap))
	}
	if snap[0].Key != e1.Key() {
		t.Errorf("unexpected remaining key: %s", snap[0].Key)
	}
}
