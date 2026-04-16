package monitor_test

import (
	"testing"
	"time"

	"github.com/deanrtaylor1/portwatch/internal/scanner"
)

// TestRateLimiter_ScanCycle_SuppressDuplicates verifies that a RateLimiter
// integrated into a scan cycle suppresses repeated alerts for the same entry.
func TestRateLimiter_ScanCycle_SuppressDuplicates(t *testing.T) {
	rl := scanner.NewRateLimiter(5 * time.Minute)

	added := []scanner.Entry{
		{LocalAddress: "0.0.0.0", LocalPort: 9090, Protocol: "tcp"},
	}

	// First pass: all entries should be allowed.
	pass1 := rl.Filter(added)
	if len(pass1) != 1 {
		t.Fatalf("pass1: expected 1, got %d", len(pass1))
	}

	// Second pass within cooldown: entries should be suppressed.
	pass2 := rl.Filter(added)
	if len(pass2) != 0 {
		t.Fatalf("pass2: expected 0, got %d", len(pass2))
	}
}

// TestRateLimiter_ZeroCooldown_AlwaysAllows verifies that a zero cooldown
// never suppresses entries.
func TestRateLimiter_ZeroCooldown_AlwaysAllows(t *testing.T) {
	rl := scanner.NewRateLimiter(0)

	e := scanner.Entry{LocalAddress: "127.0.0.1", LocalPort: 8080, Protocol: "tcp"}

	for i := 0; i < 3; i++ {
		if !rl.Allow(e.Key()) {
			t.Fatalf("iteration %d: expected allow with zero cooldown", i)
		}
	}
}
