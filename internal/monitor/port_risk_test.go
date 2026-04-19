package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iamcalledrob/portwatch/internal/scanner"
)

func makeRiskEntry(port uint16, proto scanner.Protocol) scanner.Entry {
	return scanner.Entry{Port: port, Protocol: proto, LocalAddress: "0.0.0.0"}
}

func TestPortRiskStore_Empty(t *testing.T) {
	s := NewPortRiskStore(scanner.NewClassifier(nil))
	if got := s.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty, got %d", len(got))
	}
}

func TestPortRiskStore_Record_HighPort(t *testing.T) {
	s := NewPortRiskStore(scanner.NewClassifier(nil))
	s.Record(makeRiskEntry(80, scanner.ProtocolTCP))
	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	if snap[0].Port != 80 {
		t.Errorf("expected port 80, got %d", snap[0].Port)
	}
	if snap[0].Score != 1.0 {
		t.Errorf("expected score 1.0 for well-known port, got %f", snap[0].Score)
	}
}

func TestPortRiskStore_SortedByScore(t *testing.T) {
	s := NewPortRiskStore(scanner.NewClassifier(nil))
	s.Record(makeRiskEntry(9999, scanner.ProtocolTCP))
	s.Record(makeRiskEntry(80, scanner.ProtocolTCP))
	snap := s.Snapshot()
	if snap[0].Score < snap[1].Score {
		t.Error("expected descending score order")
	}
}

func TestPortRiskStore_NilClassifier(t *testing.T) {
	s := NewPortRiskStore(nil)
	s.Record(makeRiskEntry(80, scanner.ProtocolTCP))
	if len(s.Snapshot()) != 0 {
		t.Error("expected no entries with nil classifier")
	}
}

func TestPortRiskAPI_Get(t *testing.T) {
	s := NewPortRiskStore(scanner.NewClassifier(nil))
	s.Record(makeRiskEntry(443, scanner.ProtocolTCP))
	h := NewPortRiskAPI(s)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/risk", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out []PortRiskEntry
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestPortRiskAPI_MethodNotAllowed(t *testing.T) {
	s := NewPortRiskStore(nil)
	h := NewPortRiskAPI(s)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/risk", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestPortRiskHook_RecordsAdded(t *testing.T) {
	s := NewPortRiskStore(scanner.NewClassifier(nil))
	h := NewPortRiskHook(s)
	h.OnScan([]scanner.Entry{makeRiskEntry(22, scanner.ProtocolTCP)}, nil)
	if len(s.Snapshot()) != 1 {
		t.Error("expected 1 entry after hook")
	}
}
