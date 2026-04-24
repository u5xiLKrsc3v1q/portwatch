package monitor

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// TestPortChurnHook_NilStore ensures a nil store does not panic.
func TestPortChurnHook_NilStore(t *testing.T) {
	h := NewPortChurnHook(nil)
	h.OnScan([]scanner.Entry{makeChurnEntry(1234)}, nil) // must not panic
}

// TestPortChurnStore_DefaultWindow verifies a zero window falls back to 5 minutes.
func TestPortChurnStore_DefaultWindow(t *testing.T) {
	s := NewPortChurnStore(0)
	if s.window != 5*time.Minute {
		t.Fatalf("expected default 5m window, got %v", s.window)
	}
}

// TestPortChurnIntegration_HookFeedsAPI runs a small end-to-end path:
// hook records events → store accumulates → API reflects counts.
func TestPortChurnIntegration_HookFeedsAPI(t *testing.T) {
	store := NewPortChurnStore(time.Minute)
	hook := NewPortChurnHook(store)

	e1 := makeChurnEntry(5000)
	e2 := makeChurnEntry(6000)

	// Simulate two scan cycles with churn.
	hook.OnScan([]scanner.Entry{e1}, nil)
	hook.OnScan([]scanner.Entry{e2}, []scanner.Entry{e1})
	hook.OnScan(nil, []scanner.Entry{e2})

	// e1 seen twice (added + removed), e2 seen twice (added + removed).
	if got := store.Score(e1); got != 2 {
		t.Fatalf("e1: expected 2 events, got %d", got)
	}
	if got := store.Score(e2); got != 2 {
		t.Fatalf("e2: expected 2 events, got %d", got)
	}

	snap := store.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 keys in snapshot, got %d", len(snap))
	}
}
