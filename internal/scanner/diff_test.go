package scanner

import (
	"testing"
)

func makeEntry(proto, addr string, port uint16) PortEntry {
	return PortEntry{Protocol: proto, LocalAddr: addr, Port: port}
}

func TestDiff_NoChanges(t *testing.T) {
	e := makeEntry("tcp", "00000000:0050", 80)
	prev := NewSnapshot([]PortEntry{e})
	curr := NewSnapshot([]PortEntry{e})

	result := Diff(prev, curr)
	if result.HasChanges() {
		t.Error("expected no changes, got some")
	}
}

func TestDiff_Added(t *testing.T) {
	prev := NewSnapshot([]PortEntry{
		makeEntry("tcp", "00000000:0050", 80),
	})
	curr := NewSnapshot([]PortEntry{
		makeEntry("tcp", "00000000:0050", 80),
		makeEntry("tcp", "00000000:01BB", 443),
	})

	result := Diff(prev, curr)
	if len(result.Added) != 1 {
		t.Errorf("expected 1 added entry, got %d", len(result.Added))
	}
	if result.Added[0].Port != 443 {
		t.Errorf("expected added port 443, got %d", result.Added[0].Port)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed entries, got %d", len(result.Removed))
	}
}

func TestDiff_Removed(t *testing.T) {
	prev := NewSnapshot([]PortEntry{
		makeEntry("tcp", "00000000:0050", 80),
		makeEntry("udp", "00000000:0035", 53),
	})
	curr := NewSnapshot([]PortEntry{
		makeEntry("tcp", "00000000:0050", 80),
	})

	result := Diff(prev, curr)
	if len(result.Removed) != 1 {
		t.Errorf("expected 1 removed entry, got %d", len(result.Removed))
	}
	if result.Removed[0].Port != 53 {
		t.Errorf("expected removed port 53, got %d", result.Removed[0].Port)
	}
}

func TestNewSnapshot_Key(t *testing.T) {
	e := makeEntry("tcp6", "00000000:1F90", 8080)
	if e.Key() != "tcp6|00000000:1F90" {
		t.Errorf("unexpected key: %s", e.Key())
	}
}
