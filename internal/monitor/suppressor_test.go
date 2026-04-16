package monitor

import (
	"testing"
	"time"

	"github.com/danvolchek/portwatch/internal/scanner"
)

func makeSuppEntry(port uint16) scanner.Entry {
	return scanner.Entry{Port: port, Address: "0.0.0.0", Protocol: "tcp"}
}

func TestSuppressor_NilLimiter_PassesThrough(t *testing.T) {
	s := NewSuppressor(nil)
	added := []scanner.Entry{makeSuppEntry(8080)}
	got, _ := s.Apply(added, nil)
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestSuppressor_DuplicateAdded_Suppressed(t *testing.T) {
	s := NewSuppressor(nil)
	added := []scanner.Entry{makeSuppEntry(8080)}
	s.Apply(added, nil)
	got, _ := s.Apply(added, nil)
	if len(got) != 0 {
		t.Fatalf("expected 0 after duplicate, got %d", len(got))
	}
}

func TestSuppressor_RemovedNotSuppressed(t *testing.T) {
	s := NewSuppressor(nil)
	removed := []scanner.Entry{makeSuppEntry(9090)}
	s.Apply(nil, removed)
	_, got := s.Apply(nil, removed)
	if len(got) != 1 {
		t.Fatalf("expected removed to pass, got %d", len(got))
	}
}

func TestSuppressor_WithRateLimiter_SuppressesWithinCooldown(t *testing.T) {
	limiter := scanner.NewRateLimiter(10 * time.Second)
	s := NewSuppressor(limiter)
	e := makeSuppEntry(1234)

	// first call: fingerprint is new, rate limiter allows
	got, _ := s.Apply([]scanner.Entry{e}, nil)
	if len(got) != 1 {
		t.Fatalf("expected 1 on first call, got %d", len(got))
	}

	// second call with different fingerprint but same entry: rate limiter suppresses
	// Reset guard so fingerprint doesn't suppress
	s.guard = NewFingerprintGuard()
	got, _ = s.Apply([]scanner.Entry{e}, nil)
	if len(got) != 0 {
		t.Fatalf("expected 0 within cooldown, got %d", len(got))
	}
}
