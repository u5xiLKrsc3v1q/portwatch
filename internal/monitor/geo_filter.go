package monitor

import (
	"net"
	"sync"
)

// GeoRecord holds a lightweight country/ASN tag for an IP.
type GeoRecord struct {
	IP      string
	Country string
	ASN     string
}

// GeoLookupFunc resolves an IP string to a GeoRecord.
type GeoLookupFunc func(ip string) (GeoRecord, error)

// GeoFilter annotates scan entries with geo metadata and optionally
// blocks entries whose country is in the deny list.
type GeoFilter struct {
	mu      sync.Mutex
	lookup  GeoLookupFunc
	denySet map[string]struct{}
}

// NewGeoFilter creates a GeoFilter with the provided lookup function and
// a list of ISO-3166-1 alpha-2 country codes to block.
func NewGeoFilter(fn GeoLookupFunc, deniedCountries []string) *GeoFilter {
	set := make(map[string]struct{}, len(deniedCountries))
	for _, c := range deniedCountries {
		set[c] = struct{}{}
	}
	return &GeoFilter{lookup: fn, denySet: set}
}

// IsDenied returns true when the entry's address resolves to a denied country.
func (g *GeoFilter) IsDenied(address string) bool {
	if g == nil || g.lookup == nil {
		return false
	}
	ip := net.ParseIP(address)
	if ip == nil {
		return false
	}
	rec, err := g.lookup(address)
	if err != nil {
		return false
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	_, blocked := g.denySet[rec.Country]
	return blocked
}

// Lookup exposes the raw geo record for an address.
func (g *GeoFilter) Lookup(address string) (GeoRecord, bool) {
	if g == nil || g.lookup == nil {
		return GeoRecord{}, false
	}
	rec, err := g.lookup(address)
	if err != nil {
		return GeoRecord{}, false
	}
	return rec, true
}
