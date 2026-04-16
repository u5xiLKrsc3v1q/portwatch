package scanner

import (
	"testing"
)

func TestParseProtocol_Known(t *testing.T) {
	cases := []struct {
		input    string
		expected Protocol
	}{
		{"tcp", ProtocolTCP},
		{"TCP", ProtocolTCP},
		{"tcp6", ProtocolTCP6},
		{"udp", ProtocolUDP},
		{"UDP6", ProtocolUDP6},
		{"  tcp  ", ProtocolTCP},
	}
	for _, tc := range cases {
		got := ParseProtocol(tc.input)
		if got != tc.expected {
			t.Errorf("ParseProtocol(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestParseProtocol_Unknown(t *testing.T) {
	got := ParseProtocol("sctp")
	if got != ProtocolUnknown {
		t.Errorf("expected unknown, got %q", got)
	}
}

func TestProtocol_IsValid(t *testing.T) {
	for _, p := range AllProtocols() {
		if !p.IsValid() {
			t.Errorf("expected %q to be valid", p)
		}
	}
	if ProtocolUnknown.IsValid() {
		t.Error("expected unknown to be invalid")
	}
}

func TestProtocol_IsTCP(t *testing.T) {
	if !ProtocolTCP.IsTCP() {
		t.Error("tcp should be TCP")
	}
	if !ProtocolTCP6.IsTCP() {
		t.Error("tcp6 should be TCP")
	}
	if ProtocolUDP.IsTCP() {
		t.Error("udp should not be TCP")
	}
}

func TestProtocol_IsUDP(t *testing.T) {
	if !ProtocolUDP.IsUDP() {
		t.Error("udp should be UDP")
	}
	if !ProtocolUDP6.IsUDP() {
		t.Error("udp6 should be UDP")
	}
	if ProtocolTCP.IsUDP() {
		t.Error("tcp should not be UDP")
	}
}

func TestProtocol_String(t *testing.T) {
	if ProtocolTCP.String() != "tcp" {
		t.Errorf("expected 'tcp', got %q", ProtocolTCP.String())
	}
}
