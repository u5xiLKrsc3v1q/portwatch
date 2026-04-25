package monitor

import (
	"testing"

	"github.com/ben/portwatch/internal/scanner"
)

func makeAlertCountEntry(proto, addr string, port uint16) scanner.Entry {
	return scanner.Entry{
		Protocol: scanner.Protocol(proto),
		LocalAddress: addr,
		LocalPort: port,
	}
}

func TestPortAlertCountHook_NilStore_NoOp(t *testing.T) {
	h := NewPortAlertCountHook(nil)
	// Should not panic.
	h.OnAlert(AlertEvent{
		Added: []scanner.Entry{makeAlertCountEntry("tcp", "0.0.0.0", 80)},
	})
}

func TestPortAlertCountHook_RecordsAdded(t *testing.T) {
	store := NewPortAlertCountStore()
	h := NewPortAlertCountHook(store)
	e1 := makeAlertCountEntry("tcp", "0.0.0.0", 80)
	e2 := makeAlertCountEntry("udp", "0.0.0.0", 53)
	h.OnAlert(AlertEvent{Added: []scanner.Entry{e1, e2}})
	h.OnAlert(AlertEvent{Added: []scanner.Entry{e1}})
	if got := store.Count(e1.Key()); got != 2 {
		t.Fatalf("expected 2 for e1, got %d", got)
	}
	if got := store.Count(e2.Key()); got != 1 {
		t.Fatalf("expected 1 for e2, got %d", got)
	}
}

func TestPortAlertCountHook_IgnoresRemoved(t *testing.T) {
	store := NewPortAlertCountStore()
	h := NewPortAlertCountHook(store)
	e := makeAlertCountEntry("tcp", "127.0.0.1", 9090)
	h.OnAlert(AlertEvent{Removed: []scanner.Entry{e}})
	if got := store.Count(e.Key()); got != 0 {
		t.Fatalf("expected 0 for removed entry, got %d", got)
	}
}

func TestPortAlertCountHook_EmptyEvent(t *testing.T) {
	store := NewPortAlertCountStore()
	h := NewPortAlertCountHook(store)
	h.OnAlert(AlertEvent{})
	if snap := store.Snapshot(); len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %d", len(snap))
	}
}
