package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rgst-io/portwatch/internal/scanner"
)

func makeSnapshotEntry(port uint16) scanner.Entry {
	return scanner.Entry{LocalAddress: "0.0.0.0", LocalPort: port, Protocol: scanner.TCP}
}

func TestPortSnapshotAPI_Empty(t *testing.T) {
	store := NewPortSnapshotStore()
	api := NewPortSnapshotAPI(store)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	api.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Count != 0 {
		t.Errorf("expected count 0, got %d", resp.Count)
	}
}

func TestPortSnapshotAPI_WithEntries(t *testing.T) {
	store := NewPortSnapshotStore()
	store.Update([]scanner.Entry{makeSnapshotEntry(80), makeSnapshotEntry(443)})
	api := NewPortSnapshotAPI(store)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	api.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp struct {
		Count int `json:"count"`
	}
	_ = json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Count != 2 {
		t.Errorf("expected count 2, got %d", resp.Count)
	}
}

func TestPortSnapshotAPI_MethodNotAllowed(t *testing.T) {
	store := NewPortSnapshotStore()
	api := NewPortSnapshotAPI(store)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/snapshot", nil)
	api.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestPortSnapshotHook_UpdatesStore(t *testing.T) {
	store := NewPortSnapshotStore()
	hook := NewPortSnapshotHook(store)
	entries := []scanner.Entry{makeSnapshotEntry(22)}
	hook.OnScan(entries)

	got, ts := store.Snapshot()
	if len(got) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got))
	}
	if ts.IsZero() {
		t.Error("expected non-zero updated_at")
	}
}
