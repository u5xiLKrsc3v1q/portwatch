package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeSilenceEntry(port int) scanner.Entry {
	return scanner.Entry{
		Port:     port,
		Address:  "0.0.0.0",
		Protocol: scanner.TCP,
	}
}

func TestPortSilenceStore_Empty(t *testing.T) {
	s := NewPortSilenceStore()
	if got := s.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %d", len(got))
	}
}

func TestPortSilenceStore_SilenceAndCheck(t *testing.T) {
	s := NewPortSilenceStore()
	e := makeSilenceEntry(8080)
	s.Silence(e, 0)
	if !s.IsSilenced(e) {
		t.Fatal("expected entry to be silenced")
	}
}

func TestPortSilenceStore_Unsilence(t *testing.T) {
	s := NewPortSilenceStore()
	e := makeSilenceEntry(9090)
	s.Silence(e, 0)
	s.Unsilence(e.Key())
	if s.IsSilenced(e) {
		t.Fatal("expected entry to not be silenced after unsilence")
	}
}

func TestPortSilenceStore_Expiry(t *testing.T) {
	s := NewPortSilenceStore()
	e := makeSilenceEntry(7070)
	s.Silence(e, 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if s.IsSilenced(e) {
		t.Fatal("expected silence to have expired")
	}
	if snap := s.Snapshot(); len(snap) != 0 {
		t.Fatalf("expected expired entry to be excluded from snapshot, got %d", len(snap))
	}
}

func TestPortSilenceStore_Snapshot_ExcludesExpired(t *testing.T) {
	s := NewPortSilenceStore()
	e1 := makeSilenceEntry(1111)
	e2 := makeSilenceEntry(2222)
	s.Silence(e1, 0) // no expiry
	s.Silence(e2, 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry in snapshot, got %d", len(snap))
	}
}

func TestPortSilenceAPI_Get(t *testing.T) {
	s := NewPortSilenceStore()
	s.Silence(makeSilenceEntry(3000), 0)
	api := NewPortSilenceAPI(s)
	req := httptest.NewRequest(http.MethodGet, "/silence", nil)
	rw := httptest.NewRecorder()
	api.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestPortSilenceAPI_MethodNotAllowed(t *testing.T) {
	api := NewPortSilenceAPI(NewPortSilenceStore())
	req := httptest.NewRequest(http.MethodPost, "/silence", nil)
	rw := httptest.NewRecorder()
	api.ServeHTTP(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}
