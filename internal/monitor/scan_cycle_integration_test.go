package monitor

import (
	"bytes"
	"context"
	"log"
	"testing"

	"github.com/derekg/portwatch/internal/scanner"
)

// TestScanCycle_BaselineFiltersAdded verifies that entries present in the
// baseline are not forwarded to the notifier as "added" events.
func TestScanCycle_BaselineFiltersAdded(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	state := scanner.NewStateStore()
	stub := &stubNotifier{}

	// Build a baseline that contains every entry returned by the first scan.
	// We do this by running one scan, collecting entries, then building baseline.
	entries, err := scanner.Scan(nil)
	if err != nil {
		t.Skipf("scan unavailable: %v", err)
	}

	bl := makeBaseline(entries) // helper defined in baseline_runner_test.go
	bm := NewBaselineManager(bl)

	cycle := NewScanCycle(state, nil, bm, stub, logger)
	cycle.Run(context.Background())

	// Because all existing ports are in the baseline, notifier should NOT
	// have been called (no unrecognised additions).
	if stub.sentMsg != "" {
		t.Errorf("expected no notification for baseline ports, got: %s", stub.sentMsg)
	}
}

// TestScanCycle_RunLoop_IntegrationSmoke verifies that runLoop drives
// ScanCycle.Run via a real ticker without deadlocking.
func TestScanCycle_RunLoop_IntegrationSmoke(t *testing.T) {
	var buf bytes.Buffer
	state := scanner.NewStateStore()
	cycle := NewScanCycle(state, nil, nil, nil, log.New(&buf, "", 0))

	ticker := NewTicker(10e6) // 10 ms
	ctx, cancel := context.WithTimeout(context.Background(), 40e6)
	defer cancel()

	done := make(chan struct{})
	go func() {
		runLoop(ctx, ticker, func() { cycle.Run(ctx) })
		close(done)
	}()

	select {
	case <-done:
		// loop exited cleanly after context cancellation
	case <-context.Background().Done():
		t.Fatal("runLoop did not exit")
	}
}
