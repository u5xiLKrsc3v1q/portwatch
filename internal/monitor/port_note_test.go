package monitor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPortNoteStore_SetAndGet(t *testing.T) {
	s := NewPortNoteStore()
	s.Set("tcp:8080", "dev server")
	n, ok := s.Get("tcp:8080")
	if !ok {
		t.Fatal("expected note to exist")
	}
	if n.Note != "dev server" {
		t.Errorf("expected 'dev server', got %q", n.Note)
	}
	if n.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestPortNoteStore_Get_Missing(t *testing.T) {
	s := NewPortNoteStore()
	_, ok := s.Get("tcp:9999")
	if ok {
		t.Error("expected no note for unknown key")
	}
}

func TestPortNoteStore_Delete(t *testing.T) {
	s := NewPortNoteStore()
	s.Set("tcp:443", "https")
	s.Delete("tcp:443")
	_, ok := s.Get("tcp:443")
	if ok {
		t.Error("expected note to be deleted")
	}
}

func TestPortNoteStore_Snapshot_IsCopy(t *testing.T) {
	s := NewPortNoteStore()
	s.Set("udp:53", "dns")
	snap := s.Snapshot()
	snap["udp:53"] = PortNote{Note: "mutated"}
	n, _ := s.Get("udp:53")
	if n.Note != "dns" {
		t.Error("snapshot mutation affected store")
	}
}

func TestPortNoteAPI_Get(t *testing.T) {
	s := NewPortNoteStore()
	s.Set("tcp:22", "ssh")
	handler := NewPortNoteAPI(s)

	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var result map[string]PortNote
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if result["tcp:22"].Note != "ssh" {
		t.Errorf("unexpected note: %v", result)
	}
}

func TestPortNoteAPI_Post(t *testing.T) {
	s := NewPortNoteStore()
	handler := NewPortNoteAPI(s)

	body := `{"key":"tcp:80","note":"web"}`
	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	n, ok := s.Get("tcp:80")
	if !ok || n.Note != "web" {
		t.Errorf("expected note 'web', got %+v", n)
	}
}

func TestPortNoteAPI_Post_InvalidBody(t *testing.T) {
	s := NewPortNoteStore()
	handler := NewPortNoteAPI(s)

	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader([]byte("not-json")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestPortNoteAPI_Delete(t *testing.T) {
	s := NewPortNoteStore()
	s.Set("tcp:3306", "mysql")
	handler := NewPortNoteAPI(s)

	req := httptest.NewRequest(http.MethodDelete, "/notes?key=tcp:3306", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	_, ok := s.Get("tcp:3306")
	if ok {
		t.Error("expected note to be deleted via API")
	}
}

func TestPortNoteAPI_MethodNotAllowed(t *testing.T) {
	s := NewPortNoteStore()
	handler := NewPortNoteAPI(s)

	req := httptest.NewRequest(http.MethodPatch, "/notes", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
