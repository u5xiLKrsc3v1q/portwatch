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
