package monitor

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Notifier is the interface for sending alerts.
type Notifier interface {
	Send(title, message string) error
}

// Monitor periodically scans port bindings and alerts on changes.
type Monitor struct {
	interval  time.Duration
	notifiers []Notifier
	snapshot  *scanner.Snapshot
	logger    *log.Logger
}

// New creates a new Monitor with the given poll interval and notifiers.
func New(interval time.Duration, logger *log.Logger, notifiers ...Notifier) *Monitor {
	if logger == nil {
		logger = log.Default()
	}
	return &Monitor{
		interval:  interval,
		notifiers: notifiers,
		logger:    logger,
	}
}

// Run starts the monitoring loop. It blocks until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Take initial snapshot.
	initial, err := scanner.Scan()
	if err != nil {
		return err
	}
	m.snapshot = scanner.NewSnapshot(initial)
	m.logger.Printf("portwatch: initial scan found %d listeners", len(initial))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := m.tick(); err != nil {
				m.logger.Printf("portwatch: scan error: %v", err)
			}
		}
	}
}

func (m *Monitor) tick() error {
	entries, err := scanner.Scan()
	if err != nil {
		return err
	}
	next := scanner.NewSnapshot(entries)
	diff := m.snapshot.Diff(next)
	m.snapshot = next

	for _, e := range diff.Added {
		msg := formatEntry("new listener detected", e)
		m.logger.Println(msg)
		m.notify("New Port Listener", msg)
	}
	for _, e := range diff.Removed {
		msg := formatEntry("listener removed", e)
		m.logger.Println(msg)
		m.notify("Port Listener Removed", msg)
	}
	return nil
}

func (m *Monitor) notify(title, message string) {
	for _, n := range m.notifiers {
		if err := n.Send(title, message); err != nil {
			m.logger.Printf("portwatch: notifier error: %v", err)
		}
	}
}

func formatEntry(event string, e scanner.Entry) string {
	return event + ": " + e.Protocol + " port " + e.LocalAddress
}
