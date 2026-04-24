package monitor

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

// TestPortOwnerHook_NilStore ensures a nil store does not panic.
func TestPortOwnerHook_NilStore(t *testing.T) {
	hook := NewPortOwnerHook(nil)
	hook.OnScan([]scanner.Entry{makeOwnerEntry(8080, "test")})
}

// TestPortOwnerStore_MultipleProcesses verifies multiple ports are tracked.
func TestPortOwnerStore_MultipleProcesses(t *testing.T) {
	s := NewPortOwnerStore()
	hook := NewPortOwnerHook(s)
	hook.OnScan([]scanner.Entry{
		makeOwnerEntry(80, "nginx"),
		makeOwnerEntry(443, "nginx"),
		makeOwnerEntry(5432, "postgres"),
	})
	snap := s.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(snap))
	}
	if snap[5432] != "postgres" {
		t.Fatalf("expected postgres for 5432, got %q", snap[5432])
	}
}

// TestPortOwnerStore_OverwriteProcess verifies a port's owner can change.
func TestPortOwnerStore_OverwriteProcess(t *testing.T) {
	s := NewPortOwnerStore()
	s.Record(8080, "old-process")
	s.Record(8080, "new-process")
	v, ok := s.Get(8080)
	if !ok || v != "new-process" {
		t.Fatalf("expected new-process, got %q ok=%v", v, ok)
	}
}

// TestPortOwnerHook_EmptyScan clears all owners.
func TestPortOwnerHook_EmptyScan(t *testing.T) {
	s := NewPortOwnerStore()
	hook := NewPortOwnerHook(s)
	hook.OnScan([]scanner.Entry{makeOwnerEntry(1234, "svc")})
	hook.OnScan([]scanner.Entry{})
	if len(s.Snapshot()) != 0 {
		t.Fatal("expected empty snapshot after empty scan")
	}
}
