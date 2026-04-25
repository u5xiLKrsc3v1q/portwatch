package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeLastSeenEntry(port int, addr, proto string) scanner.Entry {
	return scanner.Entry{
		LocalAddress: addr,
		LocalPort:    port,
		Protocol:     scanner.Protocol(proto),
	}
}

func TestPortLastSeenStore_Record_FirstTime(t *testing.T) {
	store := NewPortLastSeenStore()
	before := time.Now()
	store.Record([]scanner.Entry{makeLastSeenEntry(80, "0.0.0.0", "tcp")})
	after := time.Now()

	key := makeLastSeenEntry(80, "0.0.0.0", "tcp").Key()
	t2, ok := store.LastSeen(key)
	if !ok {
		t.Fatal("expected entry to be recorded")
	}
	if t2.Before(before) || t2.After(after) {
		t.Errorf("timestamp %v outside expected range [%v, %v]", t2, before, after)
	}
}

func TestPortLastSeenStore_Record_Idempotent(t *testing.T) {
	store := NewPortLastSeenStore()
	e := makeLastSeenEntry(443, "127.0.0.1", "tcp")
	store.Record([]scanner.Entry{e})
	t1, _ := store.LastSeen(e.Key())

	time.Sleep(2 * time.Millisecond)
	store.Record([]scanner.Entry{e})
	t2, _ := store.LastSeen(e.Key())

	if !t2.After(t1) {
		t.Errorf("expected second record time %v to be after first %v", t2, t1)
	}
}

func TestPortLastSeenStore_LastSeen_Missing(t *testing.T) {
	store := NewPortLastSeenStore()
	_, ok := store.LastSeen("nonexistent")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestPortLastSeenStore_Snapshot_IsCopy(t *testing.T) {
	store := NewPortLastSeenStore()
	e := makeLastSeenEntry(8080, "0.0.0.0", "tcp")
	store.Record([]scanner.Entry{e})

	snap := store.Snapshot()
	delete(snap, e.Key())

	_, ok := store.LastSeen(e.Key())
	if !ok {
		t.Error("snapshot mutation affected original store")
	}
}

func TestPortLastSeenAPI_Get(t *testing.T) {
	store := NewPortLastSeenStore()
	e := makeLastSeenEntry(22, "0.0.0.0", "tcp")
	store.Record([]scanner.Entry{e})

	handler := NewPortLastSeenAPI(store)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/last-seen", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := result[e.Key()]; !ok {
		t.Errorf("expected key %q in response", e.Key())
	}
}

func TestPortLastSeenAPI_MethodNotAllowed(t *testing.T) {
	store := NewPortLastSeenStore()
	handler := NewPortLastSeenAPI(store)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/last-seen", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
