package monitor

import (
	"testing"
	"time"
)

func TestDebouncer_Allow_FirstTime(t *testing.T) {
	d := NewDebouncer(5 * time.Second)
	if !d.Allow("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestDebouncer_Allow_WithinWindow(t *testing.T) {
	now := time.Now()
	d := NewDebouncer(5 * time.Second)
	d.now = func() time.Time { return now }
	d.Allow("key1")
	if d.Allow("key1") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestDebouncer_Allow_AfterWindow(t *testing.T) {
	now := time.Now()
	d := NewDebouncer(5 * time.Second)
	d.now = func() time.Time { return now }
	d.Allow("key1")
	d.now = func() time.Time { return now.Add(6 * time.Second) }
	if !d.Allow("key1") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestDebouncer_Allow_DifferentKeys(t *testing.T) {
	d := NewDebouncer(5 * time.Second)
	d.Allow("key1")
	if !d.Allow("key2") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestDebouncer_Expire_RemovesOldKeys(t *testing.T) {
	now := time.Now()
	d := NewDebouncer(5 * time.Second)
	d.now = func() time.Time { return now }
	d.Allow("key1")
	d.Allow("key2")
	d.now = func() time.Time { return now.Add(6 * time.Second) }
	d.Expire()
	if d.Len() != 0 {
		t.Fatalf("expected 0 keys after expire, got %d", d.Len())
	}
}

func TestDebouncer_Expire_KeepsRecentKeys(t *testing.T) {
	now := time.Now()
	d := NewDebouncer(10 * time.Second)
	d.now = func() time.Time { return now }
	d.Allow("key1")
	d.now = func() time.Time { return now.Add(3 * time.Second) }
	d.Expire()
	if d.Len() != 1 {
		t.Fatalf("expected 1 key retained, got %d", d.Len())
	}
}
