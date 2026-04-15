package monitor

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewTicker_FiresAtInterval(t *testing.T) {
	ticker := NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		// received a tick as expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("ticker did not fire within timeout")
	}
}

func TestTicker_Stop_PreventsMoreTicks(t *testing.T) {
	ticker := NewTicker(10 * time.Millisecond)
	ticker.Stop()

	// After stop, channel should not deliver ticks.
	time.Sleep(30 * time.Millisecond)
	select {
	case <-ticker.C:
		t.Fatal("received tick after Stop")
	default:
		// expected: no tick
	}
}

func TestRunLoop_CallsOnTickMultipleTimes(t *testing.T) {
	var count int64
	ticker := NewTicker(10 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Millisecond)
	defer cancel()

	go runLoop(ctx, ticker, func() {
		atomic.AddInt64(&count, 1)
	})

	<-ctx.Done()
	time.Sleep(5 * time.Millisecond) // allow final onTick to complete

	got := atomic.LoadInt64(&count)
	if got < 2 {
		t.Errorf("expected at least 2 ticks, got %d", got)
	}
}

func TestRunLoop_StopsOnContextCancel(t *testing.T) {
	var count int64
	ticker := NewTicker(5 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	runLoop(ctx, ticker, func() {
		atomic.AddInt64(&count, 1)
	})

	if atomic.LoadInt64(&count) != 0 {
		t.Errorf("expected 0 ticks after immediate cancel, got %d", count)
	}
}

func TestDefaultTickerFactory(t *testing.T) {
	ticker := DefaultTickerFactory(50 * time.Millisecond)
	if ticker == nil {
		t.Fatal("expected non-nil ticker")
	}
	ticker.Stop()
}
