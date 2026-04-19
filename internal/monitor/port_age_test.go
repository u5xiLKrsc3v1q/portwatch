package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPortAgeStore_Record_FirstTime(t *testing.T) {
	s := NewPortAgeStore()
	s.Record("tcp:8080")
	age, ok := s.Age("tcp:8080")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if age < 0 {
		t.Error("age should be non-negative")
	}
}

func TestPortAgeStore_Record_Idempotent(t *testing.T) {
	s := NewPortAgeStore()
	s.Record("tcp:9090")
	snap1 := s.Snapshot()
	time.Sleep(10 * time.Millisecond)
	s.Record("tcp:9090")
	snap2 := s.Snapshot()
	if !snap1["tcp:9090"].Equal(snap2["tcp:9090"]) {
		t.Error("second Record should not update first-seen time")
	}
}

func TestPortAgeStore_Age_Missing(t *testing.T) {
	s := NewPortAgeStore()
	_, ok := s.Age("tcp:1234")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestPortAgeStore_Snapshot_IsCopy(t *testing.T) {
	s := NewPortAgeStore()
	s.Record("udp:53")
	snap := s.Snapshot()
	delete(snap, "udp:53")
	_, ok := s.Age("udp:53")
	if !ok {
		t.Error("deleting from snapshot should not affect store")
	}
}

func TestPortAgeAPI_Get(t *testing.T) {
	s := NewPortAgeStore()
	s.Record("tcp:443")
	api := NewPortAgeAPI(s)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/port-age", nil)
	api.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var rows []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&rows); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
}

func TestPortAgeAPI_MethodNotAllowed(t *testing.T) {
	api := NewPortAgeAPI(NewPortAgeStore())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/port-age", nil)
	api.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
