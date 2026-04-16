package scanner

import (
	"fmt"
	"net"
)

// Resolver performs reverse DNS lookups on IP addresses.
type Resolver struct {
	lookup func(addr string) ([]string, error)
}

// NewResolver creates a Resolver using the standard net.LookupAddr.
func NewResolver() *Resolver {
	return &Resolver{
		lookup: net.LookupAddr,
	}
}

// NewResolverWithFunc creates a Resolver with a custom lookup function (useful for testing).
func NewResolverWithFunc(fn func(addr string) ([]string, error)) *Resolver {
	return &Resolver{lookup: fn}
}

// Resolve attempts a reverse DNS lookup for the given IP address.
// Returns the first hostname found, or the original address if none found.
func (r *Resolver) Resolve(ip string) string {
	if ip == "" || ip == "0.0.0.0" || ip == "::" {
		return ip
	}
	names, err := r.lookup(ip)
	if err != nil || len(names) == 0 {
		return ip
	}
	name := names[0]
	// Trim trailing dot from DNS names
	if len(name) > 0 && name[len(name)-1] == '.' {
		name = name[:len(name)-1]
	}
	return name
}

// ResolveEntry returns a human-readable label for an Entry's address and port.
func (r *Resolver) ResolveEntry(e Entry) string {
	host := r.Resolve(e.LocalAddress)
	return fmt.Sprintf("%s:%d", host, e.LocalPort)
}
