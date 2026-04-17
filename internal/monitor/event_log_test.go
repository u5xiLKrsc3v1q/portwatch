package monitor

import (
	"testing"
	"time"

	"github.com/sgreben/portwatch/internal/scanner"
)

func makeLogEntry(port uint16) scanner.Entry {
	return scanner.Entry{LocalPort: port, Protocol: scanner.TCP}
}

func TestEventLog_DefaultMaxSize(t *testing.T) {
	l := NewEventLog(0)
	if l.maxSize != 100 {
		t.Fatalf("expected default maxSize 100, got %d", l.maxSize)
	}
}

func TestEventLog_AppendAndLen(t *testing.T) {
	l := NewEventLog(10)
	l.Append([]scanner.Entry{makeLogEntry(80)}, nil)
	if l.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", l.Len())
	}
}

func TestEventLog_Entries_Snapshot(t *testing.T) {
	l := NewEventLog(10)
	l.Append([]scanner.Entry{makeLogEntry(443)}, nil)
	entries := l.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if entries[0].Added[0].LocalPort != 443 {
		t.Errorf("unexpected port in log entry")
	}
}

func TestEventLog_EvictsOldest(t *testing.T) {
	l := NewEventLog(3)
	for i := uint16(1); i <= 4; i++ {
		l.Append([]scanner.Entry{makeLogEntry(i)}, nil)
	}
	if l.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", l.Len())
	}
	entries := l.Entries()
	if entries[0].Added[0].LocalPort != 2 {
		t.Errorf("expected oldest evicted, first port should be 2, got %d", entries[0].Added[0].LocalPort)
	}
}

func TestEventLog_Timestamp(t *testing.T) {
	before := time.Now()
	l := NewEventLog(5)
	l.Append(nil, []scanner.Entry{makeLogEntry(22)})
	after := time.Now()
	entries := l.Entries()
	if entries[0].Time.Before(before) || entries[0].Time.After(after) {
		t.Errorf("timestamp out of expected range")
	}
}

func TestEventLog_EmptyEntries(t *testing.T) {
	l := NewEventLog(5)
	if len(l.Entries()) != 0 {
		t.Errorf("expected empty log")
	}
}
