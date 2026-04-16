package monitor

import (
	"testing"

	"github.com/danvolchek/portwatch/internal/scanner"
)

func makeFPEntry(port uint16, addr, proto string) scanner.Entry {
	return scanner.Entry{
		Port:     port,
		Address:  addr,
		Protocol: scanner.Protocol(proto),
	}
}

func TestFingerprintGuard_FirstCall_NoSuppression(t *testing.T) {
	g := NewFingerprintGuard()
	added := []scanner.Entry{makeFPEntry(80, "0.0.0.0", "tcp")}
	result := g.Filter(added, nil)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
}

func TestFingerprintGuard_SameEntries_Suppressed(t *testing.T) {
	g := NewFingerprintGuard()
	added := []scanner.Entry{makeFPEntry(80, "0.0.0.0", "tcp")}
	g.Filter(added, nil)
	result := g.Filter(added, nil)
	if len(result) != 0 {
		t.Fatalf("expected 0 entries on repeat, got %d", len(result))
	}
}

func TestFingerprintGuard_DifferentEntries_NotSuppressed(t *testing.T) {
	g := NewFingerprintGuard()
	first := []scanner.Entry{makeFPEntry(80, "0.0.0.0", "tcp")}
	second := []scanner.Entry{makeFPEntry(443, "0.0.0.0", "tcp")}
	g.Filter(first, nil)
	result := g.Filter(second, nil)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry for new fingerprint, got %d", len(result))
	}
}

func TestFingerprintGuard_RemovedAlwaysPass(t *testing.T) {
	g := NewFingerprintGuard()
	removed := []scanner.Entry{makeFPEntry(80, "0.0.0.0", "tcp")}
	g.Filter(nil, removed)
	_, got := g.Filter(nil, removed)
	if len(got) != 1 {
		t.Fatalf("expected removed entries to always pass, got %d", len(got))
	}
}

func TestFingerprintGuard_NilEntries(t *testing.T) {
	g := NewFingerprintGuard()
	added, removed := g.Filter(nil, nil)
	if len(added) != 0 || len(removed) != 0 {
		t.Fatal("expected empty results for nil input")
	}
}
