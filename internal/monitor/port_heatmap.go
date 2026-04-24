package monitor

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

// HeatmapBucket represents activity count for a time bucket.
type HeatmapBucket struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
}

// PortHeatmapStore tracks port-change activity over time buckets.
type PortHeatmapStore struct {
	mu         sync.Mutex
	bucketSize time.Duration
	maxBuckets int
	buckets    map[int64]int // unix-second bucket -> count
}

// NewPortHeatmapStore creates a store with the given bucket size and max history.
func NewPortHeatmapStore(bucketSize time.Duration, maxBuckets int) *PortHeatmapStore {
	if bucketSize <= 0 {
		bucketSize = time.Minute
	}
	if maxBuckets <= 0 {
		maxBuckets = 60
	}
	return &PortHeatmapStore{
		bucketSize: bucketSize,
		maxBuckets: maxBuckets,
		buckets:    make(map[int64]int),
	}
}

// Record increments the bucket for the given time by delta.
func (h *PortHeatmapStore) Record(t time.Time, delta int) {
	key := t.Truncate(h.bucketSize).Unix()
	h.mu.Lock()
	h.buckets[key] += delta
	h.evict()
	h.mu.Unlock()
}

// evict removes buckets older than maxBuckets * bucketSize. Must be called with lock held.
func (h *PortHeatmapStore) evict() {
	cutoff := time.Now().Add(-time.Duration(h.maxBuckets) * h.bucketSize).Unix()
	for k := range h.buckets {
		if k < cutoff {
			delete(h.buckets, k)
		}
	}
}

// Snapshot returns a sorted slice of heatmap buckets.
func (h *PortHeatmapStore) Snapshot() []HeatmapBucket {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]HeatmapBucket, 0, len(h.buckets))
	for k, v := range h.buckets {
		out = append(out, HeatmapBucket{
			Timestamp: time.Unix(k, 0).UTC(),
			Count:     v,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Timestamp.Before(out[j].Timestamp)
	})
	return out
}

// NewPortHeatmapAPI returns an HTTP handler for the heatmap store.
func NewPortHeatmapAPI(store *PortHeatmapStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(store.Snapshot())
	})
}
