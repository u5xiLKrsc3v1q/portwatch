package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRouter_PortFirstSeen_Route verifies that /api/first-seen is wired up
// when a PortFirstSeenStore is attached via NewPortFirstSeenAPI.
func TestRouter_PortFirstSeen_Route(t *testing.T) {
	store := NewPortFirstSeenStore()
	api := NewPortFirstSeenAPI(store)

	mux := http.NewServeMux()
	mux.Handle("/api/first-seen", api)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/first-seen", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /api/first-seen, got %d", rec.Code)
	}
}

// TestRouter_PortFirstSeen_UnknownRoute ensures an unregistered path returns 404.
func TestRouter_PortFirstSeen_UnknownRoute(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/api/first-seen", NewPortFirstSeenAPI(NewPortFirstSeenStore()))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown route, got %d", rec.Code)
	}
}
