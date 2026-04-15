package scanner

import (
	"testing"
	"time"
)

func TestStateStore_InitiallyEmpty(t *testing.T) {
	store := NewStateStore()
	if store.HasSnapshot() {
		t.Fatal("expected no snapshot on new store")
	}
	snap, ts := store.Get()
	if snap != nil {
		t.Fatal("expected nil snapshot")
	}
	if !ts.IsZero() {
		t.Fatal("expected zero time")
	}
}

func TestStateStore_SetAndGet(t *testing.T) {
	store := NewStateStore()
	entries := []Entry{makeEntry("tcp", "0.0.0.0", 8080, "LISTEN")}
	snap := NewSnapshot(entries)

	before := time.Now()
	store.Set(snap)
	after := time.Now()

	if !store.HasSnapshot() {
		t.Fatal("expected snapshot to be set")
	}
	got, ts := store.Get()
	if got == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if len(got.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got.Entries))
	}
	if ts.Before(before) || ts.After(after) {
		t.Errorf("updatedAt %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestStateStore_UpdateAndDiff_FirstRun(t *testing.T) {
	store := NewStateStore()
	entries := []Entry{
		makeEntry("tcp", "0.0.0.0", 80, "LISTEN"),
		makeEntry("tcp", "0.0.0.0", 443, "LISTEN"),
	}
	snap := NewSnapshot(entries)
	changes := store.UpdateAndDiff(snap)

	if len(changes.Added) != 2 {
		t.Errorf("expected 2 added on first run, got %d", len(changes.Added))
	}
	if len(changes.Removed) != 0 {
		t.Errorf("expected 0 removed on first run, got %d", len(changes.Removed))
	}
}

func TestStateStore_UpdateAndDiff_SubsequentRun(t *testing.T) {
	store := NewStateStore()

	first := NewSnapshot([]Entry{makeEntry("tcp", "0.0.0.0", 80, "LISTEN")})
	store.UpdateAndDiff(first)

	second := NewSnapshot([]Entry{makeEntry("tcp", "0.0.0.0", 443, "LISTEN")})
	changes := store.UpdateAndDiff(second)

	if len(changes.Added) != 1 || changes.Added[0].Port != 443 {
		t.Errorf("expected port 443 added, got %+v", changes.Added)
	}
	if len(changes.Removed) != 1 || changes.Removed[0].Port != 80 {
		t.Errorf("expected port 80 removed, got %+v", changes.Removed)
	}
}

func TestStateStore_ConcurrentAccess(t *testing.T) {
	store := NewStateStore()
	done := make(chan struct{})

	go func() {
		for i := 0; i < 100; i++ {
			store.Set(NewSnapshot(nil))
		}
		close(done)
	}()

	for i := 0; i < 100; i++ {
		store.Get()
	}
	<-done
}
