package monitor

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// TestPortFirstSeenStore_MultipleEntries ensures independent keys are tracked.
func TestPortFirstSeenStore_MultipleEntries(t *testing.T) {
	s := NewPortFirstSeenStore()
	ports := []int{80, 443, 8080, 9090}
	for _, p := range ports {
		s.Record(makeFirstSeenEntry(p))
	}
	snap := s.Snapshot()
	if len(snap) != len(ports) {
		t.Fatalf("expected %d entries, got %d", len(ports), len(snap))
	}
}

// TestPortFirstSeenHook_IdempotentOnRepeat verifies repeated scans do not
// overwrite the original first-seen timestamp.
func TestPortFirstSeenHook_IdempotentOnRepeat(t *testing.T) {
	s := NewPortFirstSeenStore()
	h := NewPortFirstSeenHook(s)
	e := makeFirstSeenEntry(3000)

	h.OnScan([]scanner.Entry{e})
	t1, _ := s.FirstSeen(e)

	time.Sleep(10 * time.Millisecond)
	h.OnScan([]scanner.Entry{e})
	t2, _ := s.FirstSeen(e)

	if !t1.Equal(t2) {
		t.Errorf("first-seen changed on second scan: %v -> %v", t1, t2)
	}
}

// TestPortFirstSeenStore_ProtocolDistinct ensures TCP and UDP on the same port
// are stored as separate entries.
func TestPortFirstSeenStore_ProtocolDistinct(t *testing.T) {
	s := NewPortFirstSeenStore()
	tcp := scanner.Entry{LocalAddress: "0.0.0.0", LocalPort: 53, Protocol: scanner.TCP}
	udp := scanner.Entry{LocalAddress: "0.0.0.0", LocalPort: 53, Protocol: scanner.UDP}
	s.Record(tcp)
	time.Sleep(2 * time.Millisecond)
	s.Record(udp)

	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 distinct entries, got %d", len(snap))
	}
}
