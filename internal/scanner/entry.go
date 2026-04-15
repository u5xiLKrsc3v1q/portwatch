package scanner

// Entry represents a single port binding observed on the system.
type Entry struct {
	// Protocol is "tcp", "tcp6", "udp", or "udp6".
	Protocol string

	// LocalAddress is the bound address in "host:port" form.
	LocalAddress string

	// State is the connection state string (e.g. "LISTEN").
	State string

	// PID is the owning process ID, or 0 if unknown.
	PID int

	// ProcessName is the name of the owning process, or empty if unknown.
	ProcessName string
}

// Key returns a unique string identifying this entry, used for diffing.
func (e Entry) Key() string {
	return e.Protocol + "|" + e.LocalAddress
}
