package monitor

import (
	"sync"
	"time"
)

// Debouncer suppresses rapid repeated alerts for the same key within a window.
type Debouncer struct {
	mu     sync.Mutex
	last   map[string]time.Time
	window time.Duration
	now    func() time.Time
}

// NewDebouncer creates a Debouncer with the given quiet window.
func NewDebouncer(window time.Duration) *Debouncer {
	return &Debouncer{
		last:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// Allow returns true if the key has not been seen within the window.
func (d *Debouncer) Allow(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.now()
	if t, ok := d.last[key]; ok && now.Sub(t) < d.window {
		return false
	}
	d.last[key] = now
	return true
}

// Expire removes keys whose last-seen time is older than the window.
func (d *Debouncer) Expire() {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.now()
	for k, t := range d.last {
		if now.Sub(t) >= d.window {
			delete(d.last, k)
		}
	}
}

// Len returns the number of tracked keys.
func (d *Debouncer) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.last)
}
