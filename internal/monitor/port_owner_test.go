package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeOwnerEntry(port uint16, process string) scanner.Entry {
	return scanner.Entry{Port: port, Process: process}
}

func TestPortOwnerStore_Empty(t *testing.T) {
	s := NewPortOwnerStore()
	if len(s.Snapshot()) != 0 {
		t.Fatal("expected empty snapshot")
	}
}

func TestPortOwnerStore_RecordAndGet(t *testing.T) {
	s := NewPortOwnerStore()
	s.Record(8080, "nginx")
	v, ok := s.Get(8080)
	if !ok || v != "nginx" {
		t.Fatalf("expected nginx, got %q ok=%v", v, ok)
	}
}

func TestPortOwnerStore_RecordEmpty_Removes(t *testing.T) {
	s := NewPortOwnerStore()
	s.Record(443, "caddy")
	s.Record(443, "")
	_, ok := s.Get(443)
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestPortOwnerStore_Snapshot_IsCopy(t *testing.T) {
	s := NewPortOwnerStore()
	s.Record(22, "sshd")
	snap := s.Snapshot()
	snap[22] = "other"
	v, _ := s.Get(22)
	if v != "sshd" {
		t.Fatal("snapshot mutation affected store")
	}
}

func TestPortOwnerAPI_Get(t *testing.T) {
	s := NewPortOwnerStore()
	s.Record(80, "apache")
	h := NewPortOwnerAPI(s)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out["80"] != "apache" {
		t.Fatalf("unexpected body: %v", out)
	}
}

func TestPortOwnerAPI_MethodNotAllowed(t *testing.T) {
	s := NewPortOwnerStore()
	h := NewPortOwnerAPI(s)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestPortOwnerHook_UpdatesStore(t *testing.T) {
	s := NewPortOwnerStore()
	hook := NewPortOwnerHook(s)
	hook.OnScan([]scanner.Entry{
		makeOwnerEntry(9090, "myapp"),
		makeOwnerEntry(3000, "node"),
	})
	if v, _ := s.Get(9090); v != "myapp" {
		t.Fatalf("expected myapp, got %q", v)
	}
	// Second scan removes 9090.
	hook.OnScan([]scanner.Entry{makeOwnerEntry(3000, "node")})
	if _, ok := s.Get(9090); ok {
		t.Fatal("expected 9090 to be removed after second scan")
	}
}
