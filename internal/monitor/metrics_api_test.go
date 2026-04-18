package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeMetricsAPI() (*MetricsAPI, *Metrics) {
	m := NewMetrics()
	return NewMetricsAPI(m), m
}

func TestMetricsAPI_GetMetrics_Empty(t *testing.T) {
	api, _ := makeMetricsAPI()
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var snap MetricsSnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if snap.TotalScans != 0 {
		t.Errorf("expected 0 scans, got %d", snap.TotalScans)
	}
}

func TestMetricsAPI_GetMetrics_AfterActivity(t *testing.T) {
	api, m := makeMetricsAPI()
	m.RecordScan()
	m.RecordScan()
	m.RecordAlert()

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var snap MetricsSnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if snap.TotalScans != 2 {
		t.Errorf("expected 2 scans, got %d", snap.TotalScans)
	}
	if snap.TotalAlerts != 1 {
		t.Errorf("expected 1 alert, got %d", snap.TotalAlerts)
	}
}

func TestMetricsAPI_MethodNotAllowed(t *testing.T) {
	api, _ := makeMetricsAPI()
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/metrics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
