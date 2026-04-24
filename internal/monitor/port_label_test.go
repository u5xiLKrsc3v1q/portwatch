package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jwhittle933/portwatch/internal/scanner"
)

func TestPortLabelStore_WellKnownDefaults(t *testing.T) {
	s := NewPortLabelStore()
	if got := s.Get(80); got != "HTTP" {
		t.Errorf("expected HTTP, got %q", got)
	}
	if got := s.Get(443); got != "HTTPS" {
		t.Errorf("expected HTTPS, got %q", got)
	}
}

func TestPortLabelStore_SetAndGet(t *testing.T) {
	s := NewPortLabelStore()
	s.Set(9090, "Prometheus")
	if got := s.Get(9090); got != "Prometheus" {
		t.Errorf("expected Prometheus, got %q", got)
	}
}

func TestPortLabelStore_Get_Missing(t *testing.T) {
	s := NewPortLabelStore()
	if got := s.Get(9999); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestPortLabelStore_Snapshot_IsCopy(t *testing.T) {
	s := NewPortLabelStore()
	snap := s.Snapshot()
	snap[22] = "MODIFIED"
	if s.Get(22) == "MODIFIED" {
		t.Error("snapshot modification should not affect store")
	}
}

func TestPortLabelAPI_Get(t *testing.T) {
	s := NewPortLabelStore()
	s.Set(9090, "Prometheus")
	api := NewPortLabelAPI(s)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/labels", nil)
	api.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var out map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if out["9090"] != "Prometheus" {
		t.Errorf("expected Prometheus for 9090, got %q", out["9090"])
	}
}

func TestPortLabelAPI_MethodNotAllowed(t *testing.T) {
	s := NewPortLabelStore()
	api := NewPortLabelAPI(s)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/labels", nil)
	api.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestPortLabelHook_LabelFor_Known(t *testing.T) {
	s := NewPortLabelStore()
	h := NewPortLabelHook(s)
	e := scanner.Entry{Port: 22}
	if got := h.LabelFor(e); got != "SSH" {
		t.Errorf("expected SSH, got %q", got)
	}
}

func TestPortLabelHook_LabelFor_Unknown(t *testing.T) {
	s := NewPortLabelStore()
	h := NewPortLabelHook(s)
	e := scanner.Entry{Port: 9999}
	if got := h.LabelFor(e); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestPortLabelHook_Nil_ReturnsEmpty(t *testing.T) {
	var h *PortLabelHook
	e := scanner.Entry{Port: 80}
	if got := h.LabelFor(e); got != "" {
		t.Errorf("expected empty from nil hook, got %q", got)
	}
}

func TestPortLabelHook_AnnotateAdded(t *testing.T) {
	s := NewPortLabelStore()
	h := NewPortLabelHook(s)
	entries := []scanner.Entry{{Port: 80}, {Port: 9999}, {Port: 443}}
	annotations := h.AnnotateAdded(entries)
	if annotations[80] != "HTTP" {
		t.Errorf("expected HTTP for 80")
	}
	if annotations[443] != "HTTPS" {
		t.Errorf("expected HTTPS for 443")
	}
	if _, ok := annotations[9999]; ok {
		t.Error("9999 should not be annotated")
	}
}
