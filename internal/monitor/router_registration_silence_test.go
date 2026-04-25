package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_PortSilence_Route(t *testing.T) {
	store := NewPortSilenceStore()
	store.Silence(makeSilenceEntry(5555), 0)

	mux := http.NewServeMux()
	mux.Handle("/api/silence", NewPortSilenceAPI(store))

	req := httptest.NewRequest(http.MethodGet, "/api/silence", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 from /api/silence, got %d", rw.Code)
	}
}

func TestRouter_PortSilence_UnknownRoute(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/api/silence", NewPortSilenceAPI(NewPortSilenceStore()))

	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown route, got %d", rw.Code)
	}
}
