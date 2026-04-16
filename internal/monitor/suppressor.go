package monitor

import (
	"github.com/danvolchek/portwatch/internal/scanner"
)

// Suppressor combines fingerprint-based deduplication with rate limiting
// to prevent alert fatigue from repeated identical scan cycles.
type Suppressor struct {
	guard   *FingerprintGuard
	limiter *scanner.RateLimiter
}

// NewSuppressor creates a Suppressor with the given rate limiter.
// Pass nil to skip rate limiting.
func NewSuppressor(limiter *scanner.RateLimiter) *Suppressor {
	return &Suppressor{
		guard:   NewFingerprintGuard(),
		limiter: limiter,
	}
}

// Apply filters added entries through fingerprint guard and rate limiter.
// Removed entries bypass suppression and are always returned.
func (s *Suppressor) Apply(added, removed []scanner.Entry) ([]scanner.Entry, []scanner.Entry) {
	filtered, rem := s.guard.Filter(added, removed)

	if s.limiter != nil && len(filtered) > 0 {
		filtered = s.limiter.Filter(filtered)
	}

	return filtered, rem
}
