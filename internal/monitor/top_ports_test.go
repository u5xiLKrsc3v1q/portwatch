package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTopPortsStore_Empty(t *testing.T) {
	s := NewTopPortsStore()
	if got := s.Snapshot(10); len(got) != 0 {
		t.Fatalf("expected empty, got %d", len(got))
	}
}

func TestTopPortsStore_RecordAndCount(t *testing.T) {
	s := NewTopPortsStore()
	s.Record("80", "tcp")
	s.Record("80", "tcp")
	s.Record("443", "tcp")

	snap := s.Snapshot(10)
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	if snap[0].Port != "80" || snap[0].Count != 2 {
		t.Errorf("expected port 80 with count 2, got %+v", snap[0])
	}
}

func TestTopPortsStore_Limit(t *testing.T) {
	s := NewTopPortsStore()
	for i := 0; i < 10; i++ {
		for j := 0; j <= i; j++ {
			s.Record(string(rune('a'+i)), "tcp")
		}
	}
	snap := s.Snapshot(3)
	if len(snap) != 3 {
		t.Fatalf("expected 3, got %d", len(snap))
	}
}

func TestTopPortsStore_SortedDescending(t *testing.T) {
	s := NewTopPortsStore()
	s.Record("22", "tcp")
	s.Record("8080", "tcp")
	s.Record("8080", "tcp")
	s.Record("8080", "tcp")

	snap := s.Snapshot(10)
	if snap[0].Port != "8080" {
		t.Errorf("expected 8080 first, got %s", snap[0].Port)
	}
}

func TestTopPortsAPI_Get(t *testing.T) {
	s := NewTopPortsStore()
	s.Record("80", "tcp")
	api := NewTopPortsAPI(s)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/top-ports", nil)
	api.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result []PortHitCount
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(result) != 1 || result[0].Port != "80" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestTopPortsAPI_MethodNotAllowed(t *testing.T) {
	api := NewTopPortsAPI(NewTopPortsStore())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/top-ports", nil)
	api.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
