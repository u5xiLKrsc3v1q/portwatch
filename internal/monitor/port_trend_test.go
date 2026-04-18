package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPortTrend_RecordAndGet(t *testing.T) {
	pt := NewPortTrend()
	pt.Record(8080, "tcp")
	pt.Record(8080, "tcp")
	e := pt.Get(8080, "tcp")
	if e == nil {
		t.Fatal("expected entry")
	}
	if e.SeenCount != 2 {
		t.Errorf("expected 2, got %d", e.SeenCount)
	}
}

func TestPortTrend_Get_Missing(t *testing.T) {
	pt := NewPortTrend()
	if pt.Get(9999, "udp") != nil {
		t.Error("expected nil for unseen port")
	}
}

func TestPortTrend_Snapshot_Empty(t *testing.T) {
	pt := NewPortTrend()
	if len(pt.Snapshot()) != 0 {
		t.Error("expected empty snapshot")
	}
}

func TestPortTrend_Snapshot_Multiple(t *testing.T) {
	pt := NewPortTrend()
	pt.Record(80, "tcp")
	pt.Record(443, "tcp")
	pt.Record(80, "tcp")
	snap := pt.Snapshot()
	if len(snap) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap))
	}
}

func TestPortTrendAPI_GetTrends(t *testing.T) {
	pt := NewPortTrend()
	pt.Record(80, "tcp")
	pt.Record(80, "tcp")
	pt.Record(443, "tcp")
	api := NewPortTrendAPI(pt)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trends", nil)
	api.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []PortTrendEntry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 80 {
		t.Errorf("expected port 80 first (highest count), got %d", entries[0].Port)
	}
}

func TestPortTrendAPI_MethodNotAllowed(t *testing.T) {
	api := NewPortTrendAPI(NewPortTrend())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/trends", nil)
	api.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
