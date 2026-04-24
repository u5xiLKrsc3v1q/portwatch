package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeChurnEntry(port uint16) scanner.Entry {
	return scanner.Entry{LocalAddress: "0.0.0.0", LocalPort: port, Protocol: scanner.TCP}
}

func TestPortChurnStore_Empty(t *testing.T) {
	s := NewPortChurnStore(time.Minute)
	if got := s.Score(makeChurnEntry(80)); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestPortChurnStore_RecordAndScore(t *testing.T) {
	s := NewPortChurnStore(time.Minute)
	e := makeChurnEntry(443)
	s.Record(e)
	s.Record(e)
	if got := s.Score(e); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestPortChurnStore_Eviction(t *testing.T) {
	s := NewPortChurnStore(50 * time.Millisecond)
	e := makeChurnEntry(8080)
	s.Record(e)
	time.Sleep(80 * time.Millisecond)
	if got := s.Score(e); got != 0 {
		t.Fatalf("expected 0 after eviction, got %d", got)
	}
}

func TestPortChurnStore_Snapshot(t *testing.T) {
	s := NewPortChurnStore(time.Minute)
	s.Record(makeChurnEntry(22))
	s.Record(makeChurnEntry(22))
	s.Record(makeChurnEntry(80))
	snap := s.Snapshot()
	if snap[makeChurnEntry(22).Key()] != 2 {
		t.Fatalf("expected 2 for port 22")
	}
	if snap[makeChurnEntry(80).Key()] != 1 {
		t.Fatalf("expected 1 for port 80")
	}
}

func TestPortChurnHook_OnScan(t *testing.T) {
	s := NewPortChurnStore(time.Minute)
	h := NewPortChurnHook(s)
	added := []scanner.Entry{makeChurnEntry(9000)}
	removed := []scanner.Entry{makeChurnEntry(9001)}
	h.OnScan(added, removed)
	if s.Score(makeChurnEntry(9000)) != 1 {
		t.Fatal("expected added entry to be recorded")
	}
	if s.Score(makeChurnEntry(9001)) != 1 {
		t.Fatal("expected removed entry to be recorded")
	}
}

func TestPortChurnAPI_Get(t *testing.T) {
	s := NewPortChurnStore(time.Minute)
	s.Record(makeChurnEntry(3000))
	s.Record(makeChurnEntry(3000))
	s.Record(makeChurnEntry(4000))
	api := NewPortChurnAPI(s)
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/churn", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []churnEntry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) < 2 {
		t.Fatalf("expected at least 2 entries, got %d", len(entries))
	}
	if entries[0].Count < entries[1].Count {
		t.Fatal("expected descending order")
	}
}

func TestPortChurnAPI_MethodNotAllowed(t *testing.T) {
	api := NewPortChurnAPI(NewPortChurnStore(time.Minute))
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/churn", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
