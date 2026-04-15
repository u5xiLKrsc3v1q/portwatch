package scanner

import (
	"testing"
)

func makeEntry2(port uint16, addr, proto string) Entry {
	return Entry{
		Port:         port,
		LocalAddress: addr,
		Protocol:     proto,
	}
}

func TestNewWhitelist_Nil(t *testing.T) {
	var w *Whitelist
	if w.IsAllowed(makeEntry2(80, "0.0.0.0", "tcp")) {
		t.Error("nil whitelist should not allow any entry")
	}
}

func TestWhitelist_IsAllowed_MatchPort(t *testing.T) {
	w := NewWhitelist([]WhitelistEntry{
		{Port: 22},
	})
	if !w.IsAllowed(makeEntry2(22, "0.0.0.0", "tcp")) {
		t.Error("expected port 22 to be whitelisted")
	}
}

func TestWhitelist_IsAllowed_NoMatch(t *testing.T) {
	w := NewWhitelist([]WhitelistEntry{
		{Port: 22},
	})
	if w.IsAllowed(makeEntry2(80, "0.0.0.0", "tcp")) {
		t.Error("expected port 80 to not be whitelisted")
	}
}

func TestWhitelist_IsAllowed_AddressFilter(t *testing.T) {
	w := NewWhitelist([]WhitelistEntry{
		{Port: 8080, Address: "127.0.0.1"},
	})
	if !w.IsAllowed(makeEntry2(8080, "127.0.0.1", "tcp")) {
		t.Error("expected 127.0.0.1:8080 to be whitelisted")
	}
	if w.IsAllowed(makeEntry2(8080, "0.0.0.0", "tcp")) {
		t.Error("expected 0.0.0.0:8080 to not be whitelisted")
	}
}

func TestWhitelist_IsAllowed_ProtocolFilter(t *testing.T) {
	w := NewWhitelist([]WhitelistEntry{
		{Port: 53, Protocol: "udp"},
	})
	if !w.IsAllowed(makeEntry2(53, "0.0.0.0", "udp")) {
		t.Error("expected udp:53 to be whitelisted")
	}
	if w.IsAllowed(makeEntry2(53, "0.0.0.0", "tcp")) {
		t.Error("expected tcp:53 to not be whitelisted")
	}
}

func TestWhitelistEntry_String(t *testing.T) {
	wl := WhitelistEntry{Port: 443, Address: "0.0.0.0", Protocol: "tcp"}
	s := wl.String()
	if s != "0.0.0.0:443 (tcp)" {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestWhitelistEntry_String_Defaults(t *testing.T) {
	wl := WhitelistEntry{Port: 80}
	s := wl.String()
	if s != "*:80 (any)" {
		t.Errorf("unexpected string: %s", s)
	}
}
