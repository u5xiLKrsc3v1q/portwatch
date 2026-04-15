package scanner

import "fmt"

// Protocol represents the transport-layer protocol of a listener.
type Protocol string

const (
	ProtoTCP  Protocol = "tcp"
	ProtoTCP6 Protocol = "tcp6"
	ProtoUDP  Protocol = "udp"
	ProtoUDP6 Protocol = "udp6"
)

// Entry represents a single port binding observed on the host.
type Entry struct {
	// Protocol is the transport protocol (tcp, udp, …).
	Protocol Protocol

	// Address is the local bind address (e.g. "0.0.0.0" or "127.0.0.1").
	Address string

	// Port is the local port number.
	Port uint16

	// Inode is the socket inode number from /proc/net.
	Inode uint64

	// PID is the owning process ID, populated by Enricher (0 if unknown).
	PID int

	// Process is the owning process name, populated by Enricher (empty if unknown).
	Process string
}

// Key returns a stable string that uniquely identifies this binding,
// suitable for use as a map key or snapshot identifier.
func (e Entry) Key() string {
	return fmt.Sprintf("%s:%s:%d", e.Protocol, e.Address, e.Port)
}

// String returns a human-readable representation of the entry.
func (e Entry) String() string {
	if e.Process != "" {
		return fmt.Sprintf("%s %s:%d (pid=%d %s)", e.Protocol, e.Address, e.Port, e.PID, e.Process)
	}
	return fmt.Sprintf("%s %s:%d", e.Protocol, e.Address, e.Port)
}
