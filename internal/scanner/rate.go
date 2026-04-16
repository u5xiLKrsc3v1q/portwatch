package scanner

import (
	"sync"
	"time"
)

// RateLimiter suppresses repeated alerts for the same entry within a cooldown window.
type RateLimiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
	now      func() time.Time
}

// NewRateLimiter creates a RateLimiter with the given cooldown duration.
func NewRateLimiter(cooldown time.Duration) *RateLimiter {
	return &RateLimiter{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the entry key has not been seen within the cooldown window.
// If allowed, it records the current time for that key.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	if last, ok := r.seen[key]; ok {
		if now.Sub(last) < r.cooldown {
			return false
		}
	}
	r.seen[key] = now
	return true
}

// Filter returns only the entries whose keys pass the rate limit.
func (r *RateLimiter) Filter(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if r.Allow(e.Key()) {
			out = append(out, e)
		}
	}
	return out
}

// Expire removes stale keys older than the cooldown to prevent unbounded growth.
func (r *RateLimiter) Expire() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	for k, t := range r.seen {
		if now.Sub(t) >= r.cooldown {
			delete(r.seen, k)
		}
	}
}
