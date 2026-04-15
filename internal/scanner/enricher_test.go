package scanner

import (
	"testing"
)

// buildTestMap returns a small InodePIDMap for use in tests.
func buildTestMap() InodePIDMap {
	return InodePIDMap{
		1001: 100,
		1002: 200,
	}
}

func TestEnricher_Enrich_KnownInode(t *testing.T) {
	m := buildTestMap()
	e := NewEnricherWithMap(m)

	entries := []Entry{
		{Port: 8080, Inode: 1001},
	}
	e.Enrich(entries)

	if entries[0].PID != 100 {
		t.Errorf("expected PID 100, got %d", entries[0].PID)
	}
}

func TestEnricher_Enrich_UnknownInode(t *testing.T) {
	m := buildTestMap()
	e := NewEnricherWithMap(m)

	entries := []Entry{
		{Port: 9090, Inode: 9999},
	}
	e.Enrich(entries)

	if entries[0].PID != 0 {
		t.Errorf("expected PID 0 for unknown inode, got %d", entries[0].PID)
	}
	if entries[0].Process != "" {
		t.Errorf("expected empty Process for unknown inode, got %q", entries[0].Process)
	}
}

func TestEnricher_Enrich_MultipleEntries(t *testing.T) {
	m := buildTestMap()
	e := NewEnricherWithMap(m)

	entries := []Entry{
		{Port: 80, Inode: 1001},
		{Port: 443, Inode: 1002},
		{Port: 22, Inode: 5555},
	}
	e.Enrich(entries)

	if entries[0].PID != 100 {
		t.Errorf("entry 0: expected PID 100, got %d", entries[0].PID)
	}
	if entries[1].PID != 200 {
		t.Errorf("entry 1: expected PID 200, got %d", entries[1].PID)
	}
	if entries[2].PID != 0 {
		t.Errorf("entry 2: expected PID 0, got %d", entries[2].PID)
	}
}

func TestEnricher_EnrichOne_KnownInode(t *testing.T) {
	m := buildTestMap()
	e := NewEnricherWithMap(m)

	entry := &Entry{Port: 3000, Inode: 1002}
	e.EnrichOne(entry)

	if entry.PID != 200 {
		t.Errorf("expected PID 200, got %d", entry.PID)
	}
}

func TestEnricher_EnrichOne_UnknownInode(t *testing.T) {
	m := buildTestMap()
	e := NewEnricherWithMap(m)

	entry := &Entry{Port: 3001, Inode: 0}
	e.EnrichOne(entry)

	if entry.PID != 0 {
		t.Errorf("expected PID 0, got %d", entry.PID)
	}
}

func TestNewEnricher_DoesNotPanic(t *testing.T) {
	// NewEnricher calls BuildInodePIDMap which may fail in CI; must not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("NewEnricher panicked: %v", r)
		}
	}()
	_ = NewEnricher()
}
