package monitor

import (
	"sync"
	"time"
)

// AlertThrottle suppresses repeated alerts for the same port within a time window.
type AlertThrottle struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	nowFunc func() time.Time
}

// NewAlertThrottle creates an AlertThrottle with the given suppression window.
func NewAlertThrottle(window time.Duration) *AlertThrottle {
	return &AlertThrottle{
		seen:    make(map[string]time.Time),
		window:  window,
		nowFunc: time.Now,
	}
}

// Allow returns true if the key has not been seen within the throttle window.
func (t *AlertThrottle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.nowFunc()
	if last, ok := t.seen[key]; ok && now.Sub(last) < t.window {
		return false
	}
	t.seen[key] = now
	return true
}

// Expire removes keys whose last-seen time is older than the window.
func (t *AlertThrottle) Expire() {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.nowFunc()
	for k, ts := range t.seen {
		if now.Sub(ts) >= t.window {
			delete(t.seen, k)
		}
	}
}

// Len returns the number of tracked keys.
func (t *AlertThrottle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.seen)
}
