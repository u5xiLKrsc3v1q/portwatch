package monitor

import (
	"sync"
	"time"
)

// PortTrendEntry records a port observation count over time.
type PortTrendEntry struct {
	Port      uint16
	Protocol  string
	SeenCount int
	FirstSeen time.Time
	LastSeen  time.Time
}

// PortTrend tracks how frequently ports appear across scan cycles.
type PortTrend struct {
	mu      sync.Mutex
	entries map[string]*PortTrendEntry
}

// NewPortTrend creates a new PortTrend tracker.
func NewPortTrend() *PortTrend {
	return &PortTrend{
		entries: make(map[string]*PortTrendEntry),
	}
}

// Record increments the observation count for a port+protocol key.
func (t *PortTrend) Record(port uint16, protocol string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	key := trendKey(port, protocol)
	now := time.Now()
	if e, ok := t.entries[key]; ok {
		e.SeenCount++
		e.LastSeen = now
	} else {
		t.entries[key] = &PortTrendEntry{
			Port:      port,
			Protocol:  protocol,
			SeenCount: 1,
			FirstSeen: now,
			LastSeen:  now,
		}
	}
}

// Get returns the trend entry for a port+protocol, or nil if not seen.
func (t *PortTrend) Get(port uint16, protocol string) *PortTrendEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[trendKey(port, protocol)]
	if !ok {
		return nil
	}
	copy := *e
	return &copy
}

// Snapshot returns a copy of all trend entries.
func (t *PortTrend) Snapshot() []PortTrendEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]PortTrendEntry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}

func trendKey(port uint16, protocol string) string {
	return protocol + ":" + string(rune(port))
}
