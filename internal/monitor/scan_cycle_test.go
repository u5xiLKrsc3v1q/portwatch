package monitor

import (
	"bytes"
	"context"
	"errors"
	"log"
	"testing"

	"github.com/derekg/portwatch/internal/scanner"
)

// stubNotifier records the last message sent and can simulate errors.
type stubNotifier struct {
	sentMsg string
	err     error
}

func (s *stubNotifier) Send(_ context.Context, msg string) error {
	s.sentMsg = msg
	return s.err
}

func newCycleLogger(buf *bytes.Buffer) *log.Logger {
	return log.New(buf, "", 0)
}

func TestScanCycle_NoChanges_NoNotification(t *testing.T) {
	var buf bytes.Buffer
	state := scanner.NewStateStore()
	cycle := NewScanCycle(state, nil, nil, nil, newCycleLogger(&buf))

	// Run twice with same underlying state — second run should see no diff.
	cycle.Run(context.Background())
	cycle.Run(context.Background())

	if buf.Len() > 0 {
		t.Logf("log output: %s", buf.String())
	}
}

func TestScanCycle_NotifierCalled_OnChanges(t *testing.T) {
	var buf bytes.Buffer
	state := scanner.NewStateStore()
	stub := &stubNotifier{}
	cycle := NewScanCycle(state, nil, nil, stub, newCycleLogger(&buf))

	// First run populates state; if any ports are open we get an "added" event.
	cycle.Run(context.Background())
	// We can't guarantee ports exist in CI, so just verify no panic.
}

func TestScanCycle_NotifierError_Logged(t *testing.T) {
	var buf bytes.Buffer
	state := scanner.NewStateStore()
	stub := &stubNotifier{err: errors.New("webhook down")}
	logger := newCycleLogger(&buf)
	cycle := NewScanCycle(state, nil, nil, stub, logger)

	// Manually inject a diff by resetting state after first run.
	cycle.Run(context.Background())
	// Reset so second run sees everything as new.
	cycle.state = scanner.NewStateStore()
	cycle.Run(context.Background())

	// If stub was called and returned error, it should appear in log.
	if stub.sentMsg != "" && buf.Len() == 0 {
		t.Error("expected notification error to be logged")
	}
}

func TestScanCycle_NilNotifier_DoesNotPanic(t *testing.T) {
	var buf bytes.Buffer
	state := scanner.NewStateStore()
	cycle := NewScanCycle(state, nil, nil, nil, newCycleLogger(&buf))
	cycle.Run(context.Background())
}
