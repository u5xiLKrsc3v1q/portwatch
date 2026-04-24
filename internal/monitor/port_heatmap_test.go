package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPortHeatmapStore_Empty(t *testing.T) {
	s := NewPortHeatmapStore(time.Minute, 60)
	if got := s.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %d buckets", len(got))
	}
}

func TestPortHeatmapStore_RecordAndSnapshot(t *testing.T) {
	s := NewPortHeatmapStore(time.Minute, 60)
	now := time.Now()
	s.Record(now, 3)
	s.Record(now, 2)
	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap))
	}
	if snap[0].Count != 5 {
		t.Errorf("expected count 5, got %d", snap[0].Count)
	}
}

func TestPortHeatmapStore_MultipleBuckets(t *testing.T) {
	s := NewPortHeatmapStore(time.Minute, 60)
	now := time.Now()
	s.Record(now, 1)
	s.Record(now.Add(-2*time.Minute), 4)
	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(snap))
	}
	// sorted ascending — older bucket first
	if snap[0].Count != 4 {
		t.Errorf("expected first bucket count 4, got %d", snap[0].Count)
	}
}

func TestPortHeatmapStore_Eviction(t *testing.T) {
	s := NewPortHeatmapStore(time.Minute, 2)
	now := time.Now()
	s.Record(now.Add(-10*time.Minute), 7) // too old
	s.Record(now, 1)                       // triggers eviction
	snap := s.Snapshot()
	for _, b := range snap {
		if b.Count == 7 {
			t.Error("expected old bucket to be evicted")
		}
	}
}

func TestPortHeatmapStore_DefaultParams(t *testing.T) {
	s := NewPortHeatmapStore(0, 0)
	if s.bucketSize != time.Minute {
		t.Errorf("expected default bucket size 1m, got %v", s.bucketSize)
	}
	if s.maxBuckets != 60 {
		t.Errorf("expected default maxBuckets 60, got %d", s.maxBuckets)
	}
}

func TestPortHeatmapAPI_Get(t *testing.T) {
	s := NewPortHeatmapStore(time.Minute, 60)
	s.Record(time.Now(), 3)
	h := NewPortHeatmapAPI(s)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/heatmap", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var buckets []HeatmapBucket
	if err := json.NewDecoder(rec.Body).Decode(&buckets); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(buckets) != 1 {
		t.Errorf("expected 1 bucket in response, got %d", len(buckets))
	}
}

func TestPortHeatmapAPI_MethodNotAllowed(t *testing.T) {
	s := NewPortHeatmapStore(time.Minute, 60)
	h := NewPortHeatmapAPI(s)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/heatmap", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
