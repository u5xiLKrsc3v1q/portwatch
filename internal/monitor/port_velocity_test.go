package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPortVelocity_EmptySnapshot(t *testing.T) {
	v := NewPortVelocity(time.Minute)
	snap := v.Snapshot()
	if snap.Added != 0 || snap.Removed != 0 || snap.Net != 0 {
		t.Fatalf("expected zero snapshot, got %+v", snap)
	}
}

func TestPortVelocity_RecordAndSnapshot(t *testing.T) {
	v := NewPortVelocity(time.Minute)
	v.Record(3, 1)
	snap := v.Snapshot()
	if snap.Added != 3 {
		t.Errorf("expected 3 added, got %d", snap.Added)
	}
	if snap.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", snap.Removed)
	}
	if snap.Net != 2 {
		t.Errorf("expected net 2, got %d", snap.Net)
	}
}

func TestPortVelocity_Eviction(t *testing.T) {
	v := NewPortVelocity(50 * time.Millisecond)
	v.Record(2, 0)
	time.Sleep(80 * time.Millisecond)
	v.Record(1, 0)
	snap := v.Snapshot()
	if snap.Added != 1 {
		t.Errorf("expected 1 after eviction, got %d", snap.Added)
	}
}

func TestPortVelocity_DefaultWindow(t *testing.T) {
	v := NewPortVelocity(0)
	if v.window != 5*time.Minute {
		t.Errorf("expected default 5m window, got %s", v.window)
	}
}

func TestPortVelocityAPI_Get(t *testing.T) {
	v := NewPortVelocity(time.Minute)
	v.Record(4, 2)
	api := NewPortVelocityAPI(v)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/velocity", nil)
	api.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if int(body["added"].(float64)) != 4 {
		t.Errorf("expected added=4, got %v", body["added"])
	}
}

func TestPortVelocityAPI_MethodNotAllowed(t *testing.T) {
	v := NewPortVelocity(time.Minute)
	api := NewPortVelocityAPI(v)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/velocity", nil)
	api.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
