package monitor

import (
	"testing"
	"time"

	"github.com/danvolchek/portwatch/internal/scanner"
)

func makeDebounceEntry(port uint16) scanner.Entry {
	return scanner.Entry{LocalAddress: "0.0.0.0", LocalPort: port, Protocol: scanner.TCP}
}

func TestDebounceFilter_NilPassesThrough(t *testing.T) {
	var f *DebounceFilter
	added := []scanner.Entry{makeDebounceEntry(8080)}
	gotAdded, gotRemoved := f.Filter(added, nil)
	if len(gotAdded) != 1 {
		t.Fatalf("expected 1 added, got %d", len(gotAdded))
	}
	if len(gotRemoved) != 0 {
		t.Fatalf("expected 0 removed, got %d", len(gotRemoved))
	}
}

func TestDebounceFilter_FirstCallAllowed(t *testing.T) {
	f := NewDebounceFilter(5 * time.Second)
	added := []scanner.Entry{makeDebounceEntry(9000)}
	got, _ := f.Filter(added, nil)
	if len(got) != 1 {
		t.Fatalf("expected entry to pass through on first call")
	}
}

func TestDebounceFilter_SecondCallSuppressed(t *testing.T) {
	now := time.Now()
	f := NewDebounceFilter(5 * time.Second)
	f.debouncer.now = func() time.Time { return now }
	added := []scanner.Entry{makeDebounceEntry(9001)}
	f.Filter(added, nil)
	got, _ := f.Filter(added, nil)
	if len(got) != 0 {
		t.Fatalf("expected entry to be suppressed on second call within window")
	}
}

func TestDebounceFilter_RemovedAlwaysPass(t *testing.T) {
	now := time.Now()
	f := NewDebounceFilter(5 * time.Second)
	f.debouncer.now = func() time.Time { return now }
	removed := []scanner.Entry{makeDebounceEntry(9002)}
	_, gotRemoved := f.Filter(nil, removed)
	if len(gotRemoved) != 1 {
		t.Fatalf("expected removed entry to always pass, got %d", len(gotRemoved))
	}
}
