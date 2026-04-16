package monitor

import (
	"context"
	"log"

	"github.com/danvolchek/portwatch/internal/notifier"
	"github.com/danvolchek/portwatch/internal/scanner"
)

// ScanCycle performs a single scan iteration: scan → diff → filter → notify.
type ScanCycle struct {
	scanner    func() ([]scanner.Entry, error)
	state      *scanner.StateStore
	tagFilter  *TagFilter
	suppressor *Suppressor
	notifier   notifier.Notifier
	logger     *log.Logger
}

// NewScanCycle constructs a ScanCycle with the provided dependencies.
func NewScanCycle(
	scanFn func() ([]scanner.Entry, error),
	state *scanner.StateStore,
	tagFilter *TagFilter,
	suppressor *Suppressor,
	n notifier.Notifier,
	logger *log.Logger,
) *ScanCycle {
	return &ScanCycle{
		scanner:    scanFn,
		state:      state,
		tagFilter:  tagFilter,
		suppressor: suppressor,
		notifier:   n,
		logger:     logger,
	}
}

// Run executes one scan cycle.
func (c *ScanCycle) Run(ctx context.Context) {
	entries, err := c.scanner()
	if err != nil {
		c.logger.Printf("scan error: %v", err)
		return
	}

	added, removed := c.state.UpdateAndDiff(entries)

	if c.tagFilter != nil {
		added, removed = c.tagFilter.Filter(added, removed)
	}

	if c.suppressor != nil {
		added, removed = c.suppressor.Apply(added, removed)
	}

	if len(added) == 0 && len(removed) == 0 {
		return
	}

	event := NewAlertEvent(added, removed)
	if err := c.notifier.Send(ctx, event.Summary()); err != nil {
		c.logger.Printf("notify error: %v", err)
	}
}
