package monitor

import (
	"testing"
	"time"

	"github.com/rgzr/portwatch/internal/scanner"
)

func makeThrottleEntry(port uint16) scanner.Entry {
	return scanner.Entry{LocalAddress: "0.0.0.0", LocalPort: port, Protocol: scanner.TCP}
}

func TestAlertThrottleFilter_Nil_PassesThrough(t *testing.T) {
	f := NewAlertThrottleFilter(nil)
	entries := []scanner.Entry{makeThrottleEntry(80), makeThrottleEntry(443)}
	got := f.FilterAdded(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestAlertThrottleFilter_FirstCall_Allowed(t *testing.T) {
	th := NewAlertThrottle(5 * time.Second)
	f := NewAlertThrottleFilter(th)
	entries := []scanner.Entry{makeThrottleEntry(8080)}
	got := f.FilterAdded(entries)
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestAlertThrottleFilter_SecondCall_Suppressed(t *testing.T) {
	th := NewAlertThrottle(5 * time.Second)
	f := NewAlertThrottleFilter(th)
	entries := []scanner.Entry{makeThrottleEntry(9090)}
	f.FilterAdded(entries)
	got := f.FilterAdded(entries)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestAlertThrottleFilter_Removed_AlwaysPass(t *testing.T) {
	th := NewAlertThrottle(5 * time.Second)
	f := NewAlertThrottleFilter(th)
	entries := []scanner.Entry{makeThrottleEntry(22), makeThrottleEntry(22)}
	// even duplicate removed entries pass through
	got := f.FilterRemoved(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestAlertThrottleFilter_AfterExpiry_AllowedAgain(t *testing.T) {
	// Use a very short TTL so the throttle window expires quickly.
	th := NewAlertThrottle(50 * time.Millisecond)
	f := NewAlertThrottleFilter(th)
	entries := []scanner.Entry{makeThrottleEntry(3000)}

	// First call should be allowed.
	if got := f.FilterAdded(entries); len(got) != 1 {
		t.Fatalf("first call: expected 1, got %d", len(got))
	}
	// Second call within the window should be suppressed.
	if got := f.FilterAdded(entries); len(got) != 0 {
		t.Fatalf("second call: expected 0, got %d", len(got))
	}
	// Wait for the throttle window to expire, then the entry should pass again.
	time.Sleep(100 * time.Millisecond)
	if got := f.FilterAdded(entries); len(got) != 1 {
		t.Fatalf("after expiry: expected 1, got %d", len(got))
	}
}
