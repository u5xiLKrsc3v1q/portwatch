package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/portwatch/internal/scanner"
)

func makeProtoEntry(proto scanner.Protocol) scanner.Entry {
	return scanner.Entry{Protocol: proto, LocalAddress: "0.0.0.0", LocalPort: 8080}
}

func TestPortProtoStatsStore_InitiallyZero(t *testing.T) {
	s := NewPortProtoStatsStore()
	snap := s.Snapshot()
	if snap.TCP != 0 || snap.UDP != 0 || snap.Other != 0 {
		t.Fatalf("expected zero stats, got %+v", snap)
	}
}

func TestPortProtoStatsStore_RecordTCP(t *testing.T) {
	s := NewPortProtoStatsStore()
	s.Record([]scanner.Entry{makeProtoEntry(scanner.ProtocolTCP)})
	if got := s.Snapshot().TCP; got != 1 {
		t.Fatalf("expected TCP=1, got %d", got)
	}
}

func TestPortProtoStatsStore_RecordUDP(t *testing.T) {
	s := NewPortProtoStatsStore()
	s.Record([]scanner.Entry{
		makeProtoEntry(scanner.ProtocolUDP),
		makeProtoEntry(scanner.ProtocolUDP),
	})
	if got := s.Snapshot().UDP; got != 2 {
		t.Fatalf("expected UDP=2, got %d", got)
	}
}

func TestPortProtoStatsStore_Reset(t *testing.T) {
	s := NewPortProtoStatsStore()
	s.Record([]scanner.Entry{makeProtoEntry(scanner.ProtocolTCP)})
	s.Reset()
	snap := s.Snapshot()
	if snap.TCP != 0 {
		t.Fatalf("expected zero after reset, got %+v", snap)
	}
}

func TestPortProtoStatsStore_Mixed(t *testing.T) {
	s := NewPortProtoStatsStore()
	s.Record([]scanner.Entry{
		makeProtoEntry(scanner.ProtocolTCP),
		makeProtoEntry(scanner.ProtocolUDP),
		makeProtoEntry(scanner.ProtocolTCP),
	})
	snap := s.Snapshot()
	if snap.TCP != 2 || snap.UDP != 1 {
		t.Fatalf("unexpected stats: %+v", snap)
	}
}

func TestPortProtoStatsAPI_Get(t *testing.T) {
	s := NewPortProtoStatsStore()
	s.Record([]scanner.Entry{makeProtoEntry(scanner.ProtocolTCP)})
	h := NewPortProtoStatsAPI(s)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/proto-stats", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out ProtoStats
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if out.TCP != 1 {
		t.Fatalf("expected TCP=1 in response, got %+v", out)
	}
}

func TestPortProtoStatsAPI_MethodNotAllowed(t *testing.T) {
	h := NewPortProtoStatsAPI(NewPortProtoStatsStore())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/proto-stats", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
