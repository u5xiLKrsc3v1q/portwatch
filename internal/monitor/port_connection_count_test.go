package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeCountEntry(proto scanner.Protocol, addr string, port uint16) scanner.Entry {
	return scanner.Entry{Protocol: proto, Address: addr, Port: port}
}

func TestPortConnectionCountStore_InitiallyEmpty(t *testing.T) {
	s := NewPortConnectionCountStore()
	if got := s.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(got))
	}
}

func TestPortConnectionCountStore_RecordAndCount(t *testing.T) {
	s := NewPortConnectionCountStore()
	e := makeCountEntry(scanner.TCP, "0.0.0.0", 8080)
	s.Record([]scanner.Entry{e})
	s.Record([]scanner.Entry{e})
	snap := s.Snapshot()
	if snap[e.Key()] != 2 {
		t.Fatalf("expected count 2, got %d", snap[e.Key()])
	}
}

func TestPortConnectionCountStore_Reset(t *testing.T) {
	s := NewPortConnectionCountStore()
	e := makeCountEntry(scanner.TCP, "0.0.0.0", 9090)
	s.Record([]scanner.Entry{e})
	s.Reset()
	if got := s.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty after reset, got %d", len(got))
	}
}

func TestPortConnectionCountStore_Snapshot_IsCopy(t *testing.T) {
	s := NewPortConnectionCountStore()
	e := makeCountEntry(scanner.UDP, "127.0.0.1", 5353)
	s.Record([]scanner.Entry{e})
	snap := s.Snapshot()
	snap[e.Key()] = 999
	if s.Snapshot()[e.Key()] == 999 {
		t.Fatal("snapshot mutation affected store")
	}
}

func TestPortConnectionCountAPI_Get(t *testing.T) {
	s := NewPortConnectionCountStore()
	s.Record([]scanner.Entry{makeCountEntry(scanner.TCP, "0.0.0.0", 80)})
	s.Record([]scanner.Entry{makeCountEntry(scanner.TCP, "0.0.0.0", 80)})
	s.Record([]scanner.Entry{makeCountEntry(scanner.TCP, "0.0.0.0", 443)})

	handler := NewPortConnectionCountAPI(s)
	rr := httptest.NewRecorder()
	handler(rr, httptest.NewRequest(http.MethodGet, "/", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result []PortConnectionCountEntry
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Count < result[1].Count {
		t.Fatal("expected descending sort by count")
	}
}

func TestPortConnectionCountAPI_MethodNotAllowed(t *testing.T) {
	s := NewPortConnectionCountStore()
	handler := NewPortConnectionCountAPI(s)
	rr := httptest.NewRecorder()
	handler(rr, httptest.NewRequest(http.MethodPost, "/", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestPortConnectionCountHook_NilStore(t *testing.T) {
	h := NewPortConnectionCountHook(nil)
	h.OnScan([]scanner.Entry{makeCountEntry(scanner.TCP, "0.0.0.0", 22)})
}

func TestPortConnectionCountHook_RecordsEntries(t *testing.T) {
	s := NewPortConnectionCountStore()
	h := NewPortConnectionCountHook(s)
	e := makeCountEntry(scanner.TCP, "0.0.0.0", 22)
	h.OnScan([]scanner.Entry{e})
	h.OnScan([]scanner.Entry{e})
	if s.Snapshot()[e.Key()] != 2 {
		t.Fatal("hook did not record entries correctly")
	}
}
