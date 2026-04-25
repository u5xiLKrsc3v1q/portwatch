package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeTagSummaryEntry(tags ...string) scanner.Entry {
	e := scanner.Entry{Port: 8080, Protocol: scanner.TCP}
	e.Tags = tags
	return e
}

func TestPortTagSummaryStore_InitiallyEmpty(t *testing.T) {
	s := NewPortTagSummaryStore()
	if len(s.Snapshot()) != 0 {
		t.Fatal("expected empty snapshot")
	}
}

func TestPortTagSummaryStore_RecordAndCount(t *testing.T) {
	s := NewPortTagSummaryStore()
	s.Record([]scanner.Entry{
		makeTagSummaryEntry("baseline", "whitelisted"),
		makeTagSummaryEntry("baseline"),
	})
	snap := s.Snapshot()
	if snap["baseline"] != 2 {
		t.Errorf("expected baseline=2, got %d", snap["baseline"])
	}
	if snap["whitelisted"] != 1 {
		t.Errorf("expected whitelisted=1, got %d", snap["whitelisted"])
	}
}

func TestPortTagSummaryStore_Reset(t *testing.T) {
	s := NewPortTagSummaryStore()
	s.Record([]scanner.Entry{makeTagSummaryEntry("baseline")})
	s.Reset()
	if len(s.Snapshot()) != 0 {
		t.Fatal("expected empty after reset")
	}
}

func TestPortTagSummaryStore_Snapshot_IsCopy(t *testing.T) {
	s := NewPortTagSummaryStore()
	s.Record([]scanner.Entry{makeTagSummaryEntry("x")})
	a := s.Snapshot()
	a["x"] = 999
	if s.Snapshot()["x"] != 1 {
		t.Fatal("snapshot modification affected store")
	}
}

func TestPortTagSummaryAPI_Get(t *testing.T) {
	s := NewPortTagSummaryStore()
	s.Record([]scanner.Entry{makeTagSummaryEntry("baseline")})
	handler := NewPortTagSummaryAPI(s)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result map[string]int
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result["baseline"] != 1 {
		t.Errorf("expected baseline=1, got %d", result["baseline"])
	}
}

func TestPortTagSummaryAPI_MethodNotAllowed(t *testing.T) {
	handler := NewPortTagSummaryAPI(NewPortTagSummaryStore())
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestPortTagSummaryHook_NilStore(t *testing.T) {
	h := NewPortTagSummaryHook(nil)
	h.OnScan([]scanner.Entry{makeTagSummaryEntry("x")}) // should not panic
}

func TestPortTagSummaryHook_RecordsEntries(t *testing.T) {
	s := NewPortTagSummaryStore()
	h := NewPortTagSummaryHook(s)
	h.OnScan([]scanner.Entry{makeTagSummaryEntry("whitelisted", "baseline")})
	snap := s.Snapshot()
	if snap["whitelisted"] != 1 || snap["baseline"] != 1 {
		t.Errorf("unexpected snapshot: %v", snap)
	}
}
