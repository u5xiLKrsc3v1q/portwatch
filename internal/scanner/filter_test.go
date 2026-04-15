package scanner

import (
	"testing"
)

func TestNewFilter_Empty(t *testing.T) {
	f := NewFilter(nil, nil)
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if len(f.Ports) != 0 || len(f.Addresses) != 0 {
		t.Errorf("expected empty maps, got ports=%v addrs=%v", f.Ports, f.Addresses)
	}
}

func TestFilter_Allow_BlockedPort(t *testing.T) {
	f := NewFilter([]uint16{80, 443}, nil)
	e := makeEntry("tcp", "0.0.0.0", 80, 6)
	if f.Allow(e) {
		t.Errorf("expected port 80 to be filtered out")
	}
}

func TestFilter_Allow_BlockedAddress(t *testing.T) {
	f := NewFilter(nil, []string{"127.0.0.1"})
	e := makeEntry("tcp", "127.0.0.1", 9000, 6)
	if f.Allow(e) {
		t.Errorf("expected 127.0.0.1 to be filtered out")
	}
}

func TestFilter_Allow_Permitted(t *testing.T) {
	f := NewFilter([]uint16{80}, []string{"127.0.0.1"})
	e := makeEntry("tcp", "0.0.0.0", 8080, 6)
	if !f.Allow(e) {
		t.Errorf("expected port 8080 on 0.0.0.0 to be allowed")
	}
}

func TestFilter_NilAlwaysAllows(t *testing.T) {
	var f *Filter
	e := makeEntry("tcp", "0.0.0.0", 80, 6)
	if !f.Allow(e) {
		t.Errorf("nil filter should allow everything")
	}
}

func TestFilter_ApplyToDiff(t *testing.T) {
	f := NewFilter([]uint16{22}, []string{"127.0.0.1"})

	d := Diff{
		Added: []Entry{
			makeEntry("tcp", "0.0.0.0", 22, 6),   // blocked by port
			makeEntry("tcp", "127.0.0.1", 9000, 6), // blocked by addr
			makeEntry("tcp", "0.0.0.0", 8080, 6),  // allowed
		},
		Removed: []Entry{
			makeEntry("tcp", "0.0.0.0", 22, 6), // blocked
		},
	}

	result := f.ApplyToDiff(d)

	if len(result.Added) != 1 {
		t.Errorf("expected 1 added entry after filter, got %d", len(result.Added))
	}
	if result.Added[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", result.Added[0].Port)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed entries after filter, got %d", len(result.Removed))
	}
}
