package monitor

import (
	"testing"
	"time"
)

func TestAlertHistory_DefaultMaxSize(t *testing.T) {
	h := NewAlertHistory(0)
	if h.maxSize != 100 {
		t.Fatalf("expected default max 100, got %d", h.maxSize)
	}
}

func TestAlertHistory_AppendAndLen(t *testing.T) {
	h := NewAlertHistory(10)
	h.Append(AlertRecord{Timestamp: time.Now(), Added: 1})
	h.Append(AlertRecord{Timestamp: time.Now(), Removed: 2})
	if h.Len() != 2 {
		t.Fatalf("expected 2, got %d", h.Len())
	}
}

func TestAlertHistory_EvictsOldest(t *testing.T) {
	h := NewAlertHistory(2)
	h.Append(AlertRecord{Summary: "first"})
	h.Append(AlertRecord{Summary: "second"})
	h.Append(AlertRecord{Summary: "third"})
	entries := h.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2, got %d", len(entries))
	}
	if entries[0].Summary != "second" {
		t.Errorf("expected 'second', got %q", entries[0].Summary)
	}
	if entries[1].Summary != "third" {
		t.Errorf("expected 'third', got %q", entries[1].Summary)
	}
}

func TestAlertHistory_Entries_Snapshot(t *testing.T) {
	h := NewAlertHistory(10)
	h.Append(AlertRecord{Summary: "a"})
	entries := h.Entries()
	entries[0].Summary = "mutated"
	if h.Entries()[0].Summary != "a" {
		t.Error("snapshot should not affect internal state")
	}
}

func TestAlertHistoryHook_OnAlert_NoChanges(t *testing.T) {
	h := NewAlertHistory(10)
	hook := NewAlertHistoryHook(h)
	hook.OnAlert(AlertEvent{})
	if h.Len() != 0 {
		t.Error("expected no records for empty event")
	}
}

func TestAlertHistoryHook_OnAlert_WithChanges(t *testing.T) {
	h := NewAlertHistory(10)
	hook := NewAlertHistoryHook(h)
	event := AlertEvent{Added: makeAlertEvent(1).Added}
	hook.OnAlert(event)
	if h.Len() != 1 {
		t.Fatalf("expected 1 record, got %d", h.Len())
	}
	if h.Entries()[0].Added != 1 {
		t.Errorf("expected Added=1")
	}
}
