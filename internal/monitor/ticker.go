package monitor

import (
	"context"
	"time"
)

// Ticker wraps a time.Ticker and provides a channel-based interface
// for driving periodic scan cycles in the monitor loop.
type Ticker struct {
	ticker *time.Ticker
	C      <-chan time.Time
}

// NewTicker creates a new Ticker that fires at the given interval.
func NewTicker(interval time.Duration) *Ticker {
	t := time.NewTicker(interval)
	return &Ticker{
		ticker: t,
		C:      t.C,
	}
}

// Stop halts the underlying ticker.
func (t *Ticker) Stop() {
	t.ticker.Stop()
}

// TickerFactory is a function that produces a Ticker for a given interval.
// It is used for dependency injection in tests.
type TickerFactory func(d time.Duration) *Ticker

// DefaultTickerFactory returns a real Ticker.
func DefaultTickerFactory(d time.Duration) *Ticker {
	return NewTicker(d)
}

// runLoop drives the scan-alert cycle until ctx is cancelled.
// onTick is called on every tick; it should perform a scan and send alerts.
func runLoop(ctx context.Context, ticker *Ticker, onTick func()) {
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			onTick()
		}
	}
}
