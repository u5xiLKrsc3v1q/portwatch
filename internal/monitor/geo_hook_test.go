package monitor

import (
	"testing"

	"github.com/robherley/portwatch/internal/scanner"
)

func makeGeoEntry(addr string, port uint16) scanner.Entry {
	return scanner.Entry{Address: addr, Port: port}
}

func TestGeoHook_Nil_PassesThrough(t *testing.T) {
	var h *GeoHook
	added := []scanner.Entry{makeGeoEntry("1.2.3.4", 80)}
	a, r := h.OnScan(added, nil)
	if len(a) != 1 || len(r) != 0 {
		t.Fatal("nil hook should pass through")
	}
}

func TestGeoHook_Strip_RemovesDenied(t *testing.T) {
	lookup := stubLookup(map[string]GeoRecord{
		"1.2.3.4": {IP: "1.2.3.4", Country: "CN"},
	})
	gf := NewGeoFilter(lookup, []string{"CN"})
	h := NewGeoHook(gf, true)

	added := []scanner.Entry{
		makeGeoEntry("1.2.3.4", 8080),
		makeGeoEntry("8.8.8.8", 443),
	}
	a, _ := h.OnScan(added, nil)
	if len(a) != 1 {
		t.Fatalf("expected 1 entry after strip, got %d", len(a))
	}
	if a[0].Address != "8.8.8.8" {
		t.Fatalf("unexpected entry %+v", a[0])
	}
}

func TestGeoHook_NoStrip_KeepsDenied(t *testing.T) {
	lookup := stubLookup(map[string]GeoRecord{
		"1.2.3.4": {IP: "1.2.3.4", Country: "RU"},
	})
	gf := NewGeoFilter(lookup, []string{"RU"})
	h := NewGeoHook(gf, false)

	added := []scanner.Entry{makeGeoEntry("1.2.3.4", 22)}
	a, _ := h.OnScan(added, nil)
	if len(a) != 1 {
		t.Fatal("expected entry to be kept when strip=false")
	}
}

func TestGeoHook_RemovedAlwaysPass(t *testing.T) {
	lookup := stubLookup(map[string]GeoRecord{
		"1.2.3.4": {IP: "1.2.3.4", Country: "CN"},
	})
	gf := NewGeoFilter(lookup, []string{"CN"})
	h := NewGeoHook(gf, true)

	removed := []scanner.Entry{makeGeoEntry("1.2.3.4", 80)}
	_, r := h.OnScan(nil, removed)
	if len(r) != 1 {
		t.Fatal("removed entries should always pass through")
	}
}
