package scanner

import (
	"testing"
)

func makeFingerprintEntry(proto Protocol, addr string, port uint16) Entry {
	return Entry{
		Protocol:    proto,
		LocalAddr:   addr,
		LocalPort:   port,
		RemoteAddr:  "",
		RemotePort:  0,
		State:       "LISTEN",
		Inode:       0,
		PID:         0,
		ProcessName: "",
	}
}

func TestNewFingerprint_Empty(t *testing.T) {
	fp := NewFingerprint(nil)
	if fp.Count != 0 {
		t.Errorf("expected count 0, got %d", fp.Count)
	}
	if fp.Hash == "" {
		t.Error("expected non-empty hash for empty set")
	}
}

func TestNewFingerprint_Deterministic(t *testing.T) {
	entries := []Entry{
		makeFingerprintEntry(TCP, "0.0.0.0", 80),
		makeFingerprintEntry(TCP, "0.0.0.0", 443),
	}
	fp1 := NewFingerprint(entries)
	fp2 := NewFingerprint(entries)
	if !fp1.Equal(fp2) {
		t.Error("expected identical fingerprints for same entries")
	}
}

func TestNewFingerprint_OrderIndependent(t *testing.T) {
	a := makeFingerprintEntry(TCP, "0.0.0.0", 80)
	b := makeFingerprintEntry(TCP, "0.0.0.0", 443)

	fp1 := NewFingerprint([]Entry{a, b})
	fp2 := NewFingerprint([]Entry{b, a})
	if !fp1.Equal(fp2) {
		t.Error("fingerprint should be order-independent")
	}
}

func TestNewFingerprint_DifferentEntries(t *testing.T) {
	a := NewFingerprint([]Entry{makeFingerprintEntry(TCP, "0.0.0.0", 80)})
	b := NewFingerprint([]Entry{makeFingerprintEntry(TCP, "0.0.0.0", 8080)})
	if a.Equal(b) {
		t.Error("expected different fingerprints for different entries")
	}
}

func TestFingerprint_String(t *testing.T) {
	fp := NewFingerprint([]Entry{makeFingerprintEntry(TCP, "127.0.0.1", 9000)})
	s := fp.String()
	if len(s) == 0 {
		t.Error("expected non-empty string")
	}
}

func TestFingerprint_CountMatches(t *testing.T) {
	entries := []Entry{
		makeFingerprintEntry(TCP, "0.0.0.0", 22),
		makeFingerprintEntry(UDP, "0.0.0.0", 53),
		makeFingerprintEntry(TCP, "0.0.0.0", 443),
	}
	fp := NewFingerprint(entries)
	if fp.Count != 3 {
		t.Errorf("expected count 3, got %d", fp.Count)
	}
}
