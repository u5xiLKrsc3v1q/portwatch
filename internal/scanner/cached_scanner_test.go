package scanner

import (
	"errors"
	"testing"
	"time"
)

func TestCachedScanner_CallsScanFuncOnFirstCall(t *testing.T) {
	calls := 0
	fn := func() ([]Entry, error) {
		calls++
		return []Entry{{LocalPort: 80, Protocol: TCP}}, nil
	}
	cs := NewCachedScanner(5*time.Second, fn)
	_, err := cs.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestCachedScanner_ReturnsCachedOnSecondCall(t *testing.T) {
	calls := 0
	fn := func() ([]Entry, error) {
		calls++
		return []Entry{{LocalPort: 80, Protocol: TCP}}, nil
	}
	cs := NewCachedScanner(5*time.Second, fn)
	cs.Scan()
	cs.Scan()
	if calls != 1 {
		t.Fatalf("expected 1 underlying call due to cache, got %d", calls)
	}
}

func TestCachedScanner_InvalidateForcesRescan(t *testing.T) {
	calls := 0
	fn := func() ([]Entry, error) {
		calls++
		return []Entry{{LocalPort: 443, Protocol: TCP}}, nil
	}
	cs := NewCachedScanner(5*time.Second, fn)
	cs.Scan()
	cs.Invalidate()
	cs.Scan()
	if calls != 2 {
		t.Fatalf("expected 2 calls after invalidation, got %d", calls)
	}
}

func TestCachedScanner_PropagatesScanError(t *testing.T) {
	fn := func() ([]Entry, error) {
		return nil, errors.New("scan failed")
	}
	cs := NewCachedScanner(5*time.Second, fn)
	_, err := cs.Scan()
	if err == nil {
		t.Fatal("expected error from scan func")
	}
}

func TestCachedScanner_DoesNotCacheOnError(t *testing.T) {
	calls := 0
	fn := func() ([]Entry, error) {
		calls++
		if calls == 1 {
			return nil, errors.New("transient error")
		}
		return []Entry{{LocalPort: 22, Protocol: TCP}}, nil
	}
	cs := NewCachedScanner(5*time.Second, fn)
	cs.Scan() // error, should not cache
	got, err := cs.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 entry on retry, got %d", len(got))
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}
