package scanner

import (
	"testing"
)

func TestDeduplicator_NewIsEmpty(t *testing.T) {
	d := NewDeduplicator()
	if d.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", d.Len())
	}
}

func TestDeduplicator_IsDuplicate_FirstTimeFalse(t *testing.T) {
	d := NewDeduplicator()
	e := makeEntry("0.0.0.0", 8080, "TCP")
	if d.IsDuplicate(e) {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestDeduplicator_IsDuplicate_SecondTimeTrue(t *testing.T) {
	d := NewDeduplicator()
	e := makeEntry("0.0.0.0", 8080, "TCP")
	d.IsDuplicate(e)
	if !d.IsDuplicate(e) {
		t.Fatal("second occurrence should be a duplicate")
	}
}

func TestDeduplicator_Filter_RemovesSeen(t *testing.T) {
	d := NewDeduplicator()
	entries := []Entry{
		makeEntry("0.0.0.0", 80, "TCP"),
		makeEntry("0.0.0.0", 443, "TCP"),
		makeEntry("0.0.0.0", 80, "TCP"), // duplicate within slice
	}

	out := d.Filter(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 unique entries, got %d", len(out))
	}
}

func TestDeduplicator_Filter_AllDuplicatesOnSecondCall(t *testing.T) {
	d := NewDeduplicator()
	entries := []Entry{
		makeEntry("127.0.0.1", 9090, "TCP"),
	}

	d.Filter(entries)
	out := d.Filter(entries)
	if len(out) != 0 {
		t.Fatalf("expected 0 entries on second call, got %d", len(out))
	}
}

func TestDeduplicator_Reset_ClearsState(t *testing.T) {
	d := NewDeduplicator()
	e := makeEntry("0.0.0.0", 22, "TCP")
	d.IsDuplicate(e)

	d.Reset()

	if d.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", d.Len())
	}
	if d.IsDuplicate(e) {
		t.Fatal("entry should not be duplicate after reset")
	}
}

func TestDeduplicator_Len_Tracks(t *testing.T) {
	d := NewDeduplicator()
	d.IsDuplicate(makeEntry("0.0.0.0", 1000, "TCP"))
	d.IsDuplicate(makeEntry("0.0.0.0", 2000, "TCP"))
	d.IsDuplicate(makeEntry("0.0.0.0", 1000, "TCP")) // duplicate, not counted again

	if d.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", d.Len())
	}
}
