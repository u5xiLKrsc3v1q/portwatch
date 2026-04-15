package monitor

import (
	"context"
	"log"

	"github.com/derekg/portwatch/internal/notifier"
	"github.com/derekg/portwatch/internal/scanner"
)

// ScanCycle performs a single scan-diff-alert iteration.
// It reads current port bindings, computes changes against the stored state,
// optionally filters against a baseline, and fires notifications.
type ScanCycle struct {
	state    *scanner.StateStore
	filter   *scanner.Filter
	baseline *BaselineManager
	notify   notifier.Notifier
	logger   *log.Logger
}

// NewScanCycle constructs a ScanCycle with the provided dependencies.
func NewScanCycle(
	state *scanner.StateStore,
	filter *scanner.Filter,
	baseline *BaselineManager,
	n notifier.Notifier,
	logger *log.Logger,
) *ScanCycle {
	return &ScanCycle{
		state:    state,
		filter:   filter,
		baseline: baseline,
		notify:   n,
		logger:   logger,
	}
}

// Run executes one scan cycle. It is safe to call repeatedly.
func (sc *ScanCycle) Run(ctx context.Context) {
	entries, err := scanner.Scan(sc.filter)
	if err != nil {
		sc.logger.Printf("scan error: %v", err)
		return
	}

	added, removed := sc.state.UpdateAndDiff(entries)

	if sc.baseline != nil {
		added = sc.baseline.FilterAdded(added)
	}

	event := NewAlertEvent(added, removed)
	if !event.HasChanges() {
		return
	}

	sc.logger.Printf("changes detected: %s", event.Summary())

	if sc.notify == nil {
		return
	}
	if err := sc.notify.Send(ctx, event.Summary()); err != nil {
		sc.logger.Printf("notification error: %v", err)
	}
}
