package monitor

import (
	"sync"
	"time"
)

// PortVelocity tracks how rapidly ports are appearing/disappearing.
type PortVelocity struct {
	mu      sync.Mutex
	events  []velocityEvent
	window  time.Duration
}

type velocityEvent struct {
	at    time.Time
	delta int // +1 added, -1 removed
}

// VelocitySnapshot holds a point-in-time velocity reading.
type VelocitySnapshot struct {
	Added   int
	Removed int
	Net     int
	Window  time.Duration
}

// NewPortVelocity creates a PortVelocity with the given observation window.
func NewPortVelocity(window time.Duration) *PortVelocity {
	if window <= 0 {
		window = 5 * time.Minute
	}
	return &PortVelocity{window: window}
}

// Record adds velocity events for added/removed counts.
func (v *PortVelocity) Record(added, removed int) {
	now := time.Now()
	v.mu.Lock()
	defer v.mu.Unlock()
	for i := 0; i < added; i++ {
		v.events = append(v.events, velocityEvent{at: now, delta: 1})
	}
	for i := 0; i < removed; i++ {
		v.events = append(v.events, velocityEvent{at: now, delta: -1})
	}
	v.evict(now)
}

// Snapshot returns current velocity within the window.
func (v *PortVelocity) Snapshot() VelocitySnapshot {
	now := time.Now()
	v.mu.Lock()
	defer v.mu.Unlock()
	v.evict(now)
	var added, removed int
	for _, e := range v.events {
		if e.delta > 0 {
			added++
		} else {
			removed++
		}
	}
	return VelocitySnapshot{Added: added, Removed: removed, Net: added - removed, Window: v.window}
}

func (v *PortVelocity) evict(now time.Time) {
	cutoff := now.Add(-v.window)
	i := 0
	for i < len(v.events) && v.events[i].at.Before(cutoff) {
		i++
	}
	v.events = v.events[i:]
}
