package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_PortAlertCount_Route(t *testing.T) {
	store := NewPortAlertCountStore()
	store.Record("tcp:80")

	mux := http.NewServeMux()
	mux.Handle("/api/alert-counts", NewPortAlertCountAPI(store))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/alert-counts", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRouter_PortAlertCount_UnknownRoute(t *testing.T) {
	store := NewPortAlertCountStore()

	mux := http.NewServeMux()
	mux.Handle("/api/alert-counts", NewPortAlertCountAPI(store))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
