package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeFirstSeenEntry(port int) scanner.Entry {
	return scanner.Entry{LocalAddress: "0.0.0.0", LocalPort: port, Protocol: scanner.TCP}
}

func TestPortFirstSeenStore_Record_FirstTime(t *testing.T) {
	s := NewPortFirstSeenStore()
	e := makeFirstSeenEntry(8080)
	before := time.Now()
	s.Record(e)
	after := time.Now()

	t1, ok := s.FirstSeen(e)
	if !ok {
		t.Fatal("expected entry to be recorded")
	}
	if t1.Before(before) || t1.After(after) {
		t.Errorf("first-seen time %v out of range [%v, %v]", t1, before, after)
	}
}

func TestPortFirstSeenStore_Record_Idempotent(t *testing.T) {
	s := NewPortFirstSeenStore()
	e := makeFirstSeenEntry(9090)
	s.Record(e)
	t1, _ := s.FirstSeen(e)

	time.Sleep(5 * time.Millisecond)
	s.Record(e)
	t2, _ := s.FirstSeen(e)

	if !t1.Equal(t2) {
		t.Errorf("expected idempotent: got %v then %v", t1, t2)
	}
}

func TestPortFirstSeenStore_FirstSeen_Missing(t *testing.T) {
	s := NewPortFirstSeenStore()
	_, ok := s.FirstSeen(makeFirstSeenEntry(1234))
	if ok {
		t.Fatal("expected missing entry to return false")
	}
}

func TestPortFirstSeenStore_Snapshot_IsCopy(t *testing.T) {
	s := NewPortFirstSeenStore()
	e := makeFirstSeenEntry(443)
	s.Record(e)

	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	// Mutate the copy — original must be unaffected.
	for k := range snap {
		delete(snap, k)
	}
	if len(s.Snapshot()) != 1 {
		t.Fatal("snapshot mutation affected original store")
	}
}

func TestPortFirstSeenAPI_Get(t *testing.T) {
	s := NewPortFirstSeenStore()
	s.Record(makeFirstSeenEntry(80))

	api := NewPortFirstSeenAPI(s)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/first-seen", nil)
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

func TestPortFirstSeenAPI_MethodNotAllowed(t *testing.T) {
	api := NewPortFirstSeenAPI(NewPortFirstSeenStore())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/first-seen", nil)
	api.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestPortFirstSeenHook_OnScan_Records(t *testing.T) {
	s := NewPortFirstSeenStore()
	h := NewPortFirstSeenHook(s)

	entries := []scanner.Entry{
		makeFirstSeenEntry(22),
		makeFirstSeenEntry(443),
	}
	h.OnScan(entries)

	for _, e := range entries {
		if _, ok := s.FirstSeen(e); !ok {
			t.Errorf("expected entry %v to be recorded", e)
		}
	}
}

func TestPortFirstSeenHook_NilStore(t *testing.T) {
	h := NewPortFirstSeenHook(nil)
	// Should not panic.
	h.OnScan([]scanner.Entry{makeFirstSeenEntry(8080)})
}
