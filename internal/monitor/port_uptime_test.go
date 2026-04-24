package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPortUptimeStore_Record_FirstTime(t *testing.T) {
	now := time.Now()
	s := NewPortUptimeStore()
	s.now = func() time.Time { return now }

	s.Record("tcp:0.0.0.0:8080")

	up, ok := s.Uptime("tcp:0.0.0.0:8080")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if up != 0 {
		t.Errorf("expected 0 uptime at same instant, got %v", up)
	}
}

func TestPortUptimeStore_Record_Idempotent(t *testing.T) {
	base := time.Now()
	calls := 0
	s := NewPortUptimeStore()
	s.now = func() time.Time {
		calls++
		return base.Add(time.Duration(calls-1) * time.Second)
	}

	s.Record("tcp:0.0.0.0:9090")
	s.Record("tcp:0.0.0.0:9090") // second call must not overwrite

	s.now = func() time.Time { return base.Add(5 * time.Second) }
	up, ok := s.Uptime("tcp:0.0.0.0:9090")
	if !ok {
		t.Fatal("expected key present")
	}
	if up != 5*time.Second {
		t.Errorf("expected 5s uptime, got %v", up)
	}
}

func TestPortUptimeStore_Remove(t *testing.T) {
	s := NewPortUptimeStore()
	s.Record("tcp:127.0.0.1:443")
	s.Remove("tcp:127.0.0.1:443")

	_, ok := s.Uptime("tcp:127.0.0.1:443")
	if ok {
		t.Error("expected key to be absent after remove")
	}
}

func TestPortUptimeStore_Uptime_Missing(t *testing.T) {
	s := NewPortUptimeStore()
	_, ok := s.Uptime("tcp:0.0.0.0:1234")
	if ok {
		t.Error("expected false for unknown key")
	}
}

func TestPortUptimeStore_Snapshot_IsCopy(t *testing.T) {
	base := time.Now()
	s := NewPortUptimeStore()
	s.now = func() time.Time { return base }
	s.Record("tcp:0.0.0.0:80")
	s.Record("tcp:0.0.0.0:443")

	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 records, got %d", len(snap))
	}
	// mutate the copy — original must be unaffected
	snap[0].Key = "modified"
	snap2 := s.Snapshot()
	for _, r := range snap2 {
		if r.Key == "modified" {
			t.Error("snapshot mutation affected store")
		}
	}
}

func TestPortUptimeAPI_Get(t *testing.T) {
	base := time.Now()
	s := NewPortUptimeStore()
	s.now = func() time.Time { return base }
	s.Record("tcp:0.0.0.0:8080")

	s.now = func() time.Time { return base.Add(10 * time.Second) }

	api := NewPortUptimeAPI(s)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/uptime", nil)
	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var records []UptimeRecord
	if err := json.NewDecoder(rec.Body).Decode(&records); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Seconds != 10 {
		t.Errorf("expected 10s, got %v", records[0].Seconds)
	}
}

func TestPortUptimeAPI_MethodNotAllowed(t *testing.T) {
	api := NewPortUptimeAPI(NewPortUptimeStore())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/uptime", nil)
	api.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
