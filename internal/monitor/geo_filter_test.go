package monitor

import (
	"errors"
	"testing"
)

func stubLookup(m map[string]GeoRecord) GeoLookupFunc {
	return func(ip string) (GeoRecord, error) {
		if r, ok := m[ip]; ok {
			return r, nil
		}
		return GeoRecord{}, errors.New("not found")
	}
}

func TestGeoFilter_Nil_NotDenied(t *testing.T) {
	var g *GeoFilter
	if g.IsDenied("1.2.3.4") {
		t.Fatal("nil filter should never deny")
	}
}

func TestGeoFilter_IsDenied_BlockedCountry(t *testing.T) {
	lookup := stubLookup(map[string]GeoRecord{
		"1.2.3.4": {IP: "1.2.3.4", Country: "CN"},
	})
	g := NewGeoFilter(lookup, []string{"CN", "RU"})
	if !g.IsDenied("1.2.3.4") {
		t.Fatal("expected 1.2.3.4 to be denied")
	}
}

func TestGeoFilter_IsDenied_AllowedCountry(t *testing.T) {
	lookup := stubLookup(map[string]GeoRecord{
		"8.8.8.8": {IP: "8.8.8.8", Country: "US"},
	})
	g := NewGeoFilter(lookup, []string{"CN"})
	if g.IsDenied("8.8.8.8") {
		t.Fatal("expected 8.8.8.8 to be allowed")
	}
}

func TestGeoFilter_IsDenied_InvalidIP(t *testing.T) {
	g := NewGeoFilter(stubLookup(nil), []string{"CN"})
	if g.IsDenied("not-an-ip") {
		t.Fatal("invalid IP should not be denied")
	}
}

func TestGeoFilter_IsDenied_LookupError(t *testing.T) {
	g := NewGeoFilter(stubLookup(nil), []string{"CN"})
	if g.IsDenied("1.1.1.1") {
		t.Fatal("lookup error should not deny")
	}
}

func TestGeoFilter_Lookup_Found(t *testing.T) {
	lookup := stubLookup(map[string]GeoRecord{
		"10.0.0.1": {IP: "10.0.0.1", Country: "DE", ASN: "AS1234"},
	})
	g := NewGeoFilter(lookup, nil)
	rec, ok := g.Lookup("10.0.0.1")
	if !ok || rec.Country != "DE" {
		t.Fatalf("expected DE record, got %+v ok=%v", rec, ok)
	}
}

func TestGeoFilter_Lookup_NotFound(t *testing.T) {
	g := NewGeoFilter(stubLookup(nil), nil)
	_, ok := g.Lookup("10.0.0.2")
	if ok {
		t.Fatal("expected not found")
	}
}
