package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rgst-io/portwatch/internal/scanner"
)

func makeDiffEntry(port int, proto scanner.Protocol, action string) scanner.Entry {
	return scanner.Entry{
		Port:        port,
		Protocol:    proto,
		Address:     "0.0.0.0",
		ProcessName: "testd",
	}
	_ = action // used by caller to route to added/removed
}

func TestPortStateDiffStore_DefaultMaxSize(t *testing.T) {
	s := NewPortStateDiffStore(0)
	if s.maxSize != 200 {
		t.Fatalf("expected default maxSize 200, got %d", s.maxSize)
	}
}

func TestPortStateDiffStore_RecordAdded(t *testing.T) {
	s := NewPortStateDiffStore(10)
	e := scanner.Entry{Port: 8080, Protocol: scanner.TCP, Address: "0.0.0.0"}
	s.Record([]scanner.Entry{e}, nil)
	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(snap))
	}
	if snap[0].Action != "added" {
		t.Errorf("expected action 'added', got %q", snap[0].Action)
	}
	if snap[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", snap[0].Port)
	}
}

func TestPortStateDiffStore_RecordRemoved(t *testing.T) {
	s := NewPortStateDiffStore(10)
	e := scanner.Entry{Port: 443, Protocol: scanner.TCP, Address: "127.0.0.1"}
	s.Record(nil, []scanner.Entry{e})
	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(snap))
	}
	if snap[0].Action != "removed" {
		t.Errorf("expected action 'removed', got %q", snap[0].Action)
	}
}

func TestPortStateDiffStore_Eviction(t *testing.T) {
	s := NewPortStateDiffStore(3)
	for i := 0; i < 5; i++ {
		e := scanner.Entry{Port: i + 1, Protocol: scanner.TCP}
		s.Record([]scanner.Entry{e}, nil)
	}
	snap := s.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", len(snap))
	}
	if snap[0].Port != 3 {
		t.Errorf("expected oldest surviving port 3, got %d", snap[0].Port)
	}
}

func TestPortStateDiffAPI_Get(t *testing.T) {
	s := NewPortStateDiffStore(10)
	e := scanner.Entry{Port: 9000, Protocol: scanner.UDP, Address: "::"}
	s.Record([]scanner.Entry{e}, nil)

	handler := NewPortStateDiffAPI(s)
	req := httptest.NewRequest(http.MethodGet, "/diffs", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var diffs []PortStateDiff
	if err := json.NewDecoder(rec.Body).Decode(&diffs); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(diffs) != 1 || diffs[0].Port != 9000 {
		t.Errorf("unexpected diffs: %+v", diffs)
	}
}

func TestPortStateDiffAPI_MethodNotAllowed(t *testing.T) {
	s := NewPortStateDiffStore(10)
	handler := NewPortStateDiffAPI(s)
	req := httptest.NewRequest(http.MethodPost, "/diffs", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
