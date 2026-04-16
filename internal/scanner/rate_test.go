package scanner

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow_FirstTime(t *testing.T) {
	rl := NewRateLimiter(time.Minute)
	if !rl.Allow("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestRateLimiter_Allow_WithinCooldown(t *testing.T) {
	rl := NewRateLimiter(time.Minute)
	rl.Allow("key1")
	if rl.Allow("key1") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestRateLimiter_Allow_AfterCooldown(t *testing.T) {
	now := time.Now()
	rl := NewRateLimiter(time.Minute)
	rl.now = func() time.Time { return now }
	rl.Allow("key1")

	rl.now = func() time.Time { return now.Add(2 * time.Minute) }
	if !rl.Allow("key1") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestRateLimiter_Filter(t *testing.T) {
	rl := NewRateLimiter(time.Minute)
	entries := []Entry{
		{LocalAddress: "0.0.0.0", LocalPort: 80, Protocol: "tcp"},
		{LocalAddress: "0.0.0.0", LocalPort: 443, Protocol: "tcp"},
	}

	out := rl.Filter(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}

	// second pass — all should be blocked
	out2 := rl.Filter(entries)
	if len(out2) != 0 {
		t.Fatalf("expected 0 entries on second pass, got %d", len(out2))
	}
}

func TestRateLimiter_Expire(t *testing.T) {
	now := time.Now()
	rl := NewRateLimiter(time.Minute)
	rl.now = func() time.Time { return now }
	rl.Allow("key1")

	rl.now = func() time.Time { return now.Add(2 * time.Minute) }
	rl.Expire()

	if !rl.Allow("key1") {
		t.Fatal("expected key to be allowed after expiry")
	}
}

func TestRateLimiter_IndependentKeys(t *testing.T) {
	rl := NewRateLimiter(time.Minute)
	rl.Allow("key1")
	if !rl.Allow("key2") {
		t.Fatal("expected different key to be allowed")
	}
}
