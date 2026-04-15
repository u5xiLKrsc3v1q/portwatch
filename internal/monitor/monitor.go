package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/scanner"
)

// Monitor ties together scanning, state tracking and notification.
type Monitor struct {
	cfg      *config.Config
	notifier notifier.Notifier
	scanner  *scanner.Filter
	state    *scanner.StateStore
	logger   *log.Logger
}

// New creates a Monitor from the given config and notifier.
func New(cfg *config.Config, n notifier.Notifier, logger *log.Logger) *Monitor {
	var f *scanner.Filter
	if len(cfg.BlockedPorts) > 0 || len(cfg.BlockedAddresses) > 0 {
		f = scanner.NewFilter(cfg.BlockedPorts, cfg.BlockedAddresses)
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Monitor{
		cfg:     cfg,
		notifier: n,
		scanner: f,
		state:   scanner.NewStateStore(),
		logger:  logger,
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := m.tick(); err != nil {
				m.logger.Printf("scan error: %v", err)
			}
		}
	}
}

func (m *Monitor) tick() error {
	entries, err := scanner.Scan()
	if err != nil {
		return err
	}

	if m.scanner != nil {
		filtered := entries[:0]
		for _, e := range entries {
			if m.scanner.Allow(e) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	snap := scanner.NewSnapshot(entries)
	changes := m.state.UpdateAndDiff(snap)

	for _, e := range changes.Added {
		msg := fmt.Sprintf("[portwatch] new listener: %s", formatEntry(e))
		if err := m.notifier.Send("New Port Listener", msg); err != nil {
			m.logger.Printf("notify error: %v", err)
		}
	}
	for _, e := range changes.Removed {
		m.logger.Printf("port closed: %s", formatEntry(e))
	}
	return nil
}

// formatEntry returns a human-readable string for an Entry.
func formatEntry(e scanner.Entry) string {
	return fmt.Sprintf("%s %s:%d (%s)", e.Protocol, e.Address, e.Port, e.State)
}
