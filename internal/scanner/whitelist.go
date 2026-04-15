package scanner

import "fmt"

// WhitelistEntry represents a known/expected port binding that should not trigger alerts.
type WhitelistEntry struct {
	Port    uint16
	Address string // empty means any address
	Protocol string // "tcp", "tcp6", "udp", "udp6", or empty for any
}

// Whitelist holds a set of known-good port bindings.
type Whitelist struct {
	entries []WhitelistEntry
}

// NewWhitelist creates a Whitelist from a slice of WhitelistEntry values.
func NewWhitelist(entries []WhitelistEntry) *Whitelist {
	return &Whitelist{entries: entries}
}

// IsAllowed returns true if the given Entry matches any whitelist entry.
func (w *Whitelist) IsAllowed(e Entry) bool {
	if w == nil {
		return false
	}
	for _, wl := range w.entries {
		if wl.Port != e.Port {
			continue
		}
		if wl.Address != "" && wl.Address != e.LocalAddress {
			continue
		}
		if wl.Protocol != "" && wl.Protocol != e.Protocol {
			continue
		}
		return true
	}
	return false
}

// String returns a human-readable representation of a WhitelistEntry.
func (wl WhitelistEntry) String() string {
	addr := wl.Address
	if addr == "" {
		addr = "*"
	}
	proto := wl.Protocol
	if proto == "" {
		proto = "any"
	}
	return fmt.Sprintf("%s:%d (%s)", addr, wl.Port, proto)
}
