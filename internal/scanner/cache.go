package scanner

import (
	"sync"
	"time"
)

// CacheEntry holds a cached scan result with a timestamp.
type CacheEntry struct {
	Entries   []Entry
	CachedAt  time.Time
}

// ScanCache caches the most recent scan result for a given TTL.
type ScanCache struct {
	mu      sync.RWMutex
	entry   *CacheEntry
	ttl     time.Duration
	nowFunc func() time.Time
}

// NewScanCache creates a ScanCache with the given TTL.
func NewScanCache(ttl time.Duration) *ScanCache {
	return &ScanCache{
		ttl:     ttl,
		nowFunc: time.Now,
	}
}

// Get returns cached entries if still valid, and whether the cache was hit.
func (c *ScanCache) Get() ([]Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.entry == nil {
		return nil, false
	}
	if c.nowFunc().Sub(c.entry.CachedAt) > c.ttl {
		return nil, false
	}
	return c.entry.Entries, true
}

// Set stores entries in the cache with the current timestamp.
func (c *ScanCache) Set(entries []Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry = &CacheEntry{
		Entries:  entries,
		CachedAt: c.nowFunc(),
	}
}

// Invalidate clears the cache.
func (c *ScanCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry = nil
}

// Age returns the duration since the cache was last populated, and false if
// the cache is empty.
func (c *ScanCache) Age() (time.Duration, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.entry == nil {
		return 0, false
	}
	return c.nowFunc().Sub(c.entry.CachedAt), true
}
