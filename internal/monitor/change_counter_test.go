package monitor

import (
	"testing"
	"time"
)

func TestChangeCounter_EmptyInitially(t *testing.T) {
	cc := NewChangeCounter(5 * time.Second)
	if got := cc.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestChangeCounter_RecordAndCount(t *testing.T) {
	cc := NewChangeCounter(5 * time.Second)
	cc.Record(3)
	if got := cc.Count(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestChangeCounter_EvictsOldEvents(t *testing.T) {
	base := time.Now()
	cc := NewChangeCounter(5 * time.Second)
	cc.clock = func() time.Time { return base }
	cc.Record(2)

	// advance clock beyond window
	cc.clock = func() time.Time { return base.Add(6 * time.Second) }
	if got := cc.Count(); got != 0 {
		t.Fatalf("expected 0 after eviction, got %d", got)
	}
}

func TestChangeCounter_KeepsRecentEvents(t *testing.T) {
	base := time.Now()
	cc := NewChangeCounter(10 * time.Second)
	cc.clock = func() time.Time { return base }
	cc.Record(4)

	cc.clock = func() time.Time { return base.Add(5 * time.Second) }
	cc.Record(1)

	cc.clock = func() time.Time { return base.Add(11 * time.Second) }
	// first 4 should be evicted, last 1 remains
	if got := cc.Count(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestChangeCounter_Reset(t *testing.T) {
	cc := NewChangeCounter(5 * time.Second)
	cc.Record(5)
	cc.Reset()
	if got := cc.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}
