package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/robloxrocks/portwatch/internal/scanner"
)

func makeAPILog(t *testing.T) *EventLog {
	t.Helper()
	return NewEventLog(10)
}

func TestHTTPAPI_Health(t *testing.T) {
	api := NewHTTPAPI(":0", makeAPILog(t))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	api.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %q", body["status"])
	}
}

func TestHTTPAPI_Events_Empty(t *testing.T) {
	api := NewHTTPAPI(":0", makeAPILog(t))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	api.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["count"].(float64) != 0 {
		t.Errorf("expected count 0")
	}
}

func TestHTTPAPI_Events_WithEntries(t *testing.T) {
	log := makeAPILog(t)
	log.Append(AlertEvent{
		Timestamp: time.Now(),
		Added: []scanner.Entry{{Port: 8080, Protocol: scanner.TCP}},
	})
	api := NewHTTPAPI(":0", log)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	api.mux.ServeHTTP(rec, req)
	var body map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["count"].(float64) != 1 {
		t.Errorf("expected count 1, got %v", body["count"])
	}
}

func TestHTTPAPI_StartStop(t *testing.T) {
	api := NewHTTPAPI("127.0.0.1:0", makeAPILog(t))
	if err := api.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := api.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}
