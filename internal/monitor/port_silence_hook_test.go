package monitor

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestPortSilenceHook_NilStore_PassesThrough(t *testing.T) {
	h := NewPortSilenceHook(nil)
	entries := []scanner.Entry{makeSilenceEntry(80), makeSilenceEntry(443)}
	got := h.FilterAdded(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestPortSilenceHook_RemovesSilenced(t *testing.T) {
	store := NewPortSilenceStore()
	e := makeSilenceEntry(8080)
	store.Silence(e, 0)
	h := NewPortSilenceHook(store)
	entries := []scanner.Entry{e, makeSilenceEntry(9090)}
	got := h.FilterAdded(entries)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry after filtering, got %d", len(got))
	}
	if got[0].Port != 9090 {
		t.Fatalf("expected port 9090 to remain, got %d", got[0].Port)
	}
}

func TestPortSilenceHook_AllowsUnsilenced(t *testing.T) {
	store := NewPortSilenceStore()
	h := NewPortSilenceHook(store)
	entries := []scanner.Entry{makeSilenceEntry(22), makeSilenceEntry(80)}
	got := h.FilterAdded(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestPortSilenceHook_EmptyInput(t *testing.T) {
	store := NewPortSilenceStore()
	h := NewPortSilenceHook(store)
	got := h.FilterAdded([]scanner.Entry{})
	if len(got) != 0 {
		t.Fatalf("expected empty, got %d", len(got))
	}
}
