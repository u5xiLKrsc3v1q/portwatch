package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func makeScanHistoryAPI(t *testing.T) (*ScanHistoryAPI, *ScanHistory, *http.ServeMux) {
	t.Helper()
	h := NewScanHistory(10)
	mux := http.NewServeMux()
	api := NewScanHistoryAPI(h, mux)
	return api, h, mux
}

func TestScanHistoryAPI_Empty(t *testing.T) {
	_, _, mux := makeScanHistoryAPI(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["count"].(float64) != 0 {
		t.Errorf("expected count 0, got %v", resp["count"])
	}
}

func TestScanHistoryAPI_WithRecords(t *testing.T) {
	_, h, mux := makeScanHistoryAPI(t)
	h.Append(ScanRecord{At: time.Now(), PortCount: 3, Changed: true})
	h.Append(ScanRecord{At: time.Now(), PortCount: 5, Changed: false})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["count"].(float64) != 2 {
		t.Errorf("expected count 2, got %v", resp["count"])
	}
}

func TestScanHistoryAPI_MethodNotAllowed(t *testing.T) {
	_, _, mux := makeScanHistoryAPI(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/history", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
