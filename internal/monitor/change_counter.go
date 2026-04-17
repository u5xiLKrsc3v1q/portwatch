package monitor

import (
	"sync"
	"time"
)

// ChangeCounter tracks how many add/remove events have occurred
// within a rolling time window, useful for burst detection.
type ChangeCounter struct {
	mu      sync.Mutex
	window  time.Duration
	events  []time.Time
	clock   func() time.Time
}

// NewChangeCounter creates a ChangeCounter with the given rolling window.
func NewChangeCounter(window time.Duration) *ChangeCounter {
	return &ChangeCounter{
		window: window,
		clock:  time.Now,
	}
}

// Record registers one or more change events at the current time.
func (c *ChangeCounter) Record(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	for i := 0; i < n; i++ {
		c.events = append(c.events, now)
	}
	c.evict(now)
}

// Count returns the number of events within the rolling window.
func (c *ChangeCounter) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict(c.clock())
	return len(c.events)
}

// Reset clears all recorded events.
func (c *ChangeCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = c.events[:0]
}

// evict removes events older than the window. Must be called with lock held.
func (c *ChangeCounter) evict(now time.Time) {
	cutoff := now.Add(-c.window)
	i := 0
	for i < len(c.events) && c.events[i].Before(cutoff) {
		i++
	}
	c.events = c.events[i:]
}
