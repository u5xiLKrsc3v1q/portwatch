package scanner

import (
	"testing"
	"time"
)

func makeCacheEntries() []Entry {
	return []Entry{
		{LocalAddress: "0.0.0.0", LocalPort: 8080, Protocol: TCP},
		{LocalAddress: "127.0.0.1", LocalPort: 9090, Protocol: UDP},
	}
}

func TestScanCache_MissWhenEmpty(t *testing.T) {
	c := NewScanCache(5 * time.Second)
	_, hit := c.Get()
	if hit {
		t.Fatal("expected cache miss on empty cache")
	}
}

func TestScanCache_HitAfterSet(t *testing.T) {
	c := NewScanCache(5 * time.Second)
	entries := makeCacheEntries()
	c.Set(entries)
	got, hit := c.Get()
	if !hit {
		t.Fatal("expected cache hit")
	}
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
}

func TestScanCache_MissAfterTTLExpired(t *testing.T) {
	now := time.Now()
	c := NewScanCache(1 * time.Second)
	c.nowFunc = func() time.Time { return now }
	c.Set(makeCacheEntries())

	// Advance time past TTL
	c.nowFunc = func() time.Time { return now.Add(2 * time.Second) }
	_, hit := c.Get()
	if hit {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestScanCache_HitWithinTTL(t *testing.T) {
	now := time.Now()
	c := NewScanCache(10 * time.Second)
	c.nowFunc = func() time.Time { return now }
	c.Set(makeCacheEntries())

	c.nowFunc = func() time.Time { return now.Add(5 * time.Second) }
	_, hit := c.Get()
	if !hit {
		t.Fatal("expected cache hit within TTL")
	}
}

func TestScanCache_Invalidate(t *testing.T) {
	c := NewScanCache(10 * time.Second)
	c.Set(makeCacheEntries())
	c.Invalidate()
	_, hit := c.Get()
	if hit {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestScanCache_SetOverwritesPrevious(t *testing.T) {
	c := NewScanCache(10 * time.Second)
	c.Set(makeCacheEntries())
	newEntries := []Entry{{LocalAddress: "::1", LocalPort: 443, Protocol: TCP}}
	c.Set(newEntries)
	got, hit := c.Get()
	if !hit {
		t.Fatal("expected cache hit")
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 entry after overwrite, got %d", len(got))
	}
}
