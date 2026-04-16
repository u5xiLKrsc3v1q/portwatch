package scanner

import "time"

// ScanFunc is a function that performs a raw scan and returns entries.
type ScanFunc func() ([]Entry, error)

// CachedScanner wraps a ScanFunc with a ScanCache to avoid redundant scans.
type CachedScanner struct {
	cache    *ScanCache
	scanFunc ScanFunc
}

// NewCachedScanner creates a CachedScanner with the given TTL and scan function.
func NewCachedScanner(ttl time.Duration, fn ScanFunc) *CachedScanner {
	return &CachedScanner{
		cache:    NewScanCache(ttl),
		scanFunc: fn,
	}
}

// Scan returns cached entries if valid, otherwise invokes the underlying scan function.
func (cs *CachedScanner) Scan() ([]Entry, error) {
	if entries, hit := cs.cache.Get(); hit {
		return entries, nil
	}
	entries, err := cs.scanFunc()
	if err != nil {
		return nil, err
	}
	cs.cache.Set(entries)
	return entries, nil
}

// Invalidate forces the next Scan call to bypass the cache.
func (cs *CachedScanner) Invalidate() {
	cs.cache.Invalidate()
}
