package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/robalyx/portwatch/internal/scanner"
)

func makeAlertHistoryAPI(t *testing.T) (*AlertHistory, *AlertHistoryAPI) {
	t.Helper()
	h := NewAlertHistory(50)
	return h, NewAlertHistoryAPI(h)
}

func TestAlertHistoryAPI_Empty(t *testing.T) {
	_, api := makeAlertHistoryAPI(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	api.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var out []AlertRecord
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty, got %d entries", len(out))
	}
}

func TestAlertHistoryAPI_WithRecords(t *testing.T) {
	h, api := makeAlertHistoryAPI(t)
	e := scanner.Entry{Port: 9090, Address: "0.0.0.0"}
	h.Append(AlertRecord{Timestamp: time.Now(), Added: []scanner.Entry{e}})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	api.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var out []AlertRecord
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 record, got %d", len(out))
	}
}

func TestAlertHistoryAPI_MethodNotAllowed(t *testing.T) {
	_, api := makeAlertHistoryAPI(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/alerts", nil)
	api.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
