package monitor

import (
	"testing"
	"time"
)

func TestAlertThrottle_Allow_FirstTime(t *testing.T) {
	th := NewAlertThrottle(5 * time.Second)
	if !th.Allow("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAlertThrottle_Allow_WithinWindow(t *testing.T) {
	now := time.Now()
	th := NewAlertThrottle(5 * time.Second)
	th.nowFunc = func() time.Time { return now }
	th.Allow("key1")
	if th.Allow("key1") {
		t.Fatal("expected second call within window to be suppressed")
	}
}

func TestAlertThrottle_Allow_AfterWindow(t *testing.T) {
	now := time.Now()
	th := NewAlertThrottle(5 * time.Second)
	th.nowFunc = func() time.Time { return now }
	th.Allow("key1")
	th.nowFunc = func() time.Time { return now.Add(6 * time.Second) }
	if !th.Allow("key1") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestAlertThrottle_Allow_DifferentKeys(t *testing.T) {
	th := NewAlertThrottle(5 * time.Second)
	th.Allow("key1")
	if !th.Allow("key2") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestAlertThrottle_Expire_RemovesOld(t *testing.T) {
	now := time.Now()
	th := NewAlertThrottle(5 * time.Second)
	th.nowFunc = func() time.Time { return now }
	th.Allow("key1")
	th.Allow("key2")
	th.nowFunc = func() time.Time { return now.Add(6 * time.Second) }
	th.Expire()
	if th.Len() != 0 {
		t.Fatalf("expected 0 keys after expire, got %d", th.Len())
	}
}

func TestAlertThrottle_Expire_KeepsRecent(t *testing.T) {
	now := time.Now()
	th := NewAlertThrottle(10 * time.Second)
	th.nowFunc = func() time.Time { return now }
	th.Allow("key1")
	th.nowFunc = func() time.Time { return now.Add(3 * time.Second) }
	th.Allow("key2")
	th.nowFunc = func() time.Time { return now.Add(11 * time.Second) }
	th.Expire()
	if th.Len() != 1 {
		t.Fatalf("expected 1 key remaining, got %d", th.Len())
	}
}
