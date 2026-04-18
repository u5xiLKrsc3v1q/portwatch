package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeProcessEntry(port uint16, pid int, proc string) scanner.Entry {
	return scanner.Entry{
		Port:        port,
		Address:     "0.0.0.0",
		PID:         pid,
		ProcessName: proc,
		Protocol:    scanner.TCP,
	}
}

func TestProcessMapAPI_Empty(t *testing.T) {
	store := NewProcessMapStore()
	api := NewProcessMapAPI(store)

	rr := httptest.NewRecorder()
	api.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/processes", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var out []scanner.Entry
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty, got %d entries", len(out))
	}
}

func TestProcessMapAPI_WithEntries(t *testing.T) {
	store := NewProcessMapStore()
	store.Update([]scanner.Entry{
		makeProcessEntry(8080, 123, "nginx"),
		makeProcessEntry(443, 456, "caddy"),
	})
	api := NewProcessMapAPI(store)

	rr := httptest.NewRecorder()
	api.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/processes", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var out []scanner.Entry
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestProcessMapAPI_MethodNotAllowed(t *testing.T) {
	store := NewProcessMapStore()
	api := NewProcessMapAPI(store)

	rr := httptest.NewRecorder()
	api.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/processes", nil))

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestProcessMapHook_UpdatesStore(t *testing.T) {
	store := NewProcessMapStore()
	hook := NewProcessMapHook(store)

	hook.OnScan([]scanner.Entry{makeProcessEntry(9090, 999, "myapp")})

	snap := store.Snapshot()
	if len(snap) != 1 {
		t.Errorf("expected 1 entry, got %d", len(snap))
	}
	if snap[0].ProcessName != "myapp" {
		t.Errorf("unexpected process name: %s", snap[0].ProcessName)
	}
}
