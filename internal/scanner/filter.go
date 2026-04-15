package scanner

import "strings"

// Filter holds criteria for excluding port entries from alerts.
type Filter struct {
	// Ports is a set of port numbers to ignore (e.g. well-known services).
	Ports map[uint16]struct{}
	// Addresses is a set of local addresses to ignore (e.g. "127.0.0.1").
	Addresses map[string]struct{}
}

// NewFilter constructs a Filter from human-readable slices.
// ports is a list of port numbers; addresses is a list of IP strings.
func NewFilter(ports []uint16, addresses []string) *Filter {
	f := &Filter{
		Ports:     make(map[uint16]struct{}, len(ports)),
		Addresses: make(map[string]struct{}, len(addresses)),
	}
	for _, p := range ports {
		f.Ports[p] = struct{}{}
	}
	for _, a := range addresses {
		f.Addresses[strings.TrimSpace(a)] = struct{}{}
	}
	return f
}

// Allow returns true when the entry should be reported (i.e. not filtered out).
func (f *Filter) Allow(e Entry) bool {
	if f == nil {
		return true
	}
	if _, blocked := f.Ports[e.Port]; blocked {
		return false
	}
	if _, blocked := f.Addresses[e.LocalAddr]; blocked {
		return false
	}
	return true
}

// ApplyToDiff removes filtered entries from a Diff in place and returns it.
func (f *Filter) ApplyToDiff(d Diff) Diff {
	if f == nil {
		return d
	}
	filtered := Diff{
		Added:   make([]Entry, 0, len(d.Added)),
		Removed: make([]Entry, 0, len(d.Removed)),
	}
	for _, e := range d.Added {
		if f.Allow(e) {
			filtered.Added = append(filtered.Added, e)
		}
	}
	for _, e := range d.Removed {
		if f.Allow(e) {
			filtered.Removed = append(filtered.Removed, e)
		}
	}
	return filtered
}
