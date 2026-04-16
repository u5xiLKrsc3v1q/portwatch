package scanner

import "strings"

// Protocol represents a network protocol type.
type Protocol string

const (
	ProtocolTCP  Protocol = "tcp"
	ProtocolTCP6 Protocol = "tcp6"
	ProtocolUDP  Protocol = "udp"
	ProtocolUDP6 Protocol = "udp6"
	ProtocolUnknown Protocol = "unknown"
)

// AllProtocols returns all supported protocol values.
func AllProtocols() []Protocol {
	return []Protocol{ProtocolTCP, ProtocolTCP6, ProtocolUDP, ProtocolUDP6}
}

// ParseProtocol parses a string into a Protocol, case-insensitively.
func ParseProtocol(s string) Protocol {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "tcp":
		return ProtocolTCP
	case "tcp6":
		return ProtocolTCP6
	case "udp":
		return ProtocolUDP
	case "udp6":
		return ProtocolUDP6
	default:
		return ProtocolUnknown
	}
}

// IsValid reports whether the protocol is a known supported value.
func (p Protocol) IsValid() bool {
	switch p {
	case ProtocolTCP, ProtocolTCP6, ProtocolUDP, ProtocolUDP6:
		return true
	}
	return false
}

// IsTCP reports whether the protocol is TCP or TCP6.
func (p Protocol) IsTCP() bool {
	return p == ProtocolTCP || p == ProtocolTCP6
}

// IsUDP reports whether the protocol is UDP or UDP6.
func (p Protocol) IsUDP() bool {
	return p == ProtocolUDP || p == ProtocolUDP6
}

// String returns the string representation of the protocol.
func (p Protocol) String() string {
	return string(p)
}
