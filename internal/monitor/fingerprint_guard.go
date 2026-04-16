package monitor

import (
	"log"
	"sync"

	"github.com/derekg/portwatch/internal/scanner"
)

// FingerprintGuard tracks the fingerprint of the last scan result
// and suppresses notifications when nothing has changed.
type FingerprintGuard struct {
	mu      sync.Mutex
	last    scanner.Fingerprint
	hasLast bool
	logger  *log.Logger
}

// NewFingerprintGuard creates a new guard with no prior fingerprint.
func NewFingerprintGuard(logger *log.Logger) *FingerprintGuard {
	if logger == nil {
		logger = log.Default()
	}
	return &FingerprintGuard{logger: logger}
}

// Changed returns true if the entries differ from the last recorded fingerprint.
// It also updates the stored fingerprint when a change is detected.
func (g *FingerprintGuard) Changed(entries []scanner.Entry) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	current := scanner.NewFingerprint(entries)

	if !g.hasLast {
		g.last = current
		g.hasLast = true
		g.logger.Printf("[fingerprint] initial snapshot: %s", current)
		return true
	}

	if g.last.Equal(current) {
		return false
	}

	g.logger.Printf("[fingerprint] changed: %s -> %s", g.last, current)
	g.last = current
	return true
}

// Current returns the most recently stored fingerprint.
func (g *FingerprintGuard) Current() (scanner.Fingerprint, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.last, g.hasLast
}
