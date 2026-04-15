package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewBaseline_MissingFile(t *testing.T) {
	b, err := NewBaseline(filepath.Join(t.TempDir(), "baseline.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.Entries()) != 0 {
		t.Errorf("expected empty baseline, got %d entries", len(b.Entries()))
	}
}

func TestNewBaseline_InvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "baseline.json")
	_ = os.WriteFile(p, []byte("not json"), 0o644)
	_, err := NewBaseline(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestBaseline_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "baseline.json")

	entries := []Entry{
		{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 8080, State: "LISTEN"},
		{Protocol: "tcp6", LocalAddress: "::", LocalPort: 443, State: "LISTEN"},
	}

	b, _ := NewBaseline(p)
	if err := b.Save(entries); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Reload from disk
	b2, err := NewBaseline(p)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if len(b2.Entries()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(b2.Entries()))
	}
}

func TestBaseline_Contains(t *testing.T) {
	p := filepath.Join(t.TempDir(), "baseline.json")
	b, _ := NewBaseline(p)

	e := Entry{Protocol: "tcp", LocalAddress: "127.0.0.1", LocalPort: 9090, State: "LISTEN"}
	_ = b.Save([]Entry{e})

	if !b.Contains(e) {
		t.Error("expected Contains to return true")
	}
	other := Entry{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 22, State: "LISTEN"}
	if b.Contains(other) {
		t.Error("expected Contains to return false for unknown entry")
	}
}

func TestBaseline_FileContents(t *testing.T) {
	p := filepath.Join(t.TempDir(), "baseline.json")
	b, _ := NewBaseline(p)
	entries := []Entry{{Protocol: "udp", LocalAddress: "0.0.0.0", LocalPort: 53, State: "UNCONN"}}
	_ = b.Save(entries)

	data, _ := os.ReadFile(p)
	var loaded []Entry
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("file not valid JSON: %v", err)
	}
	if len(loaded) != 1 || loaded[0].LocalPort != 53 {
		t.Errorf("unexpected file contents: %+v", loaded)
	}
}
