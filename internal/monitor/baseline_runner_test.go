package monitor

import (
	"path/filepath"
	"testing"

	"github.com/iamcalledned/portwatch/internal/scanner"
)

func makeBaseline(t *testing.T, entries []scanner.Entry) *scanner.Baseline {
	t.Helper()
	p := filepath.Join(t.TempDir(), "baseline.json")
	b, err := scanner.NewBaseline(p)
	if err != nil {
		t.Fatalf("NewBaseline: %v", err)
	}
	if len(entries) > 0 {
		if err := b.Save(entries); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}
	return b
}

func TestBaselineManager_IsBaseline_True(t *testing.T) {
	e := scanner.Entry{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 80, State: "LISTEN"}
	b := makeBaseline(t, []scanner.Entry{e})
	bm := NewBaselineManager(b, false)
	if !bm.IsBaseline(e) {
		t.Error("expected IsBaseline true")
	}
}

func TestBaselineManager_IsBaseline_False(t *testing.T) {
	bm := NewBaselineManager(makeBaseline(t, nil), false)
	e := scanner.Entry{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 9999, State: "LISTEN"}
	if bm.IsBaseline(e) {
		t.Error("expected IsBaseline false for unknown entry")
	}
}

func TestBaselineManager_NilBaseline(t *testing.T) {
	bm := NewBaselineManager(nil, false)
	e := scanner.Entry{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 22, State: "LISTEN"}
	if bm.IsBaseline(e) {
		t.Error("nil baseline should never match")
	}
}

func TestBaselineManager_FilterAdded(t *testing.T) {
	existing := scanner.Entry{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 22, State: "LISTEN"}
	newEntry := scanner.Entry{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 4444, State: "LISTEN"}
	b := makeBaseline(t, []scanner.Entry{existing})
	bm := NewBaselineManager(b, false)

	result := bm.FilterAdded([]scanner.Entry{existing, newEntry})
	if len(result) != 1 || result[0].LocalPort != 4444 {
		t.Errorf("expected only new entry, got %+v", result)
	}
}

func TestBaselineManager_SaveIfNeeded_Once(t *testing.T) {
	p := filepath.Join(t.TempDir(), "baseline.json")
	b, _ := scanner.NewBaseline(p)
	bm := NewBaselineManager(b, true)

	entries := []scanner.Entry{
		{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 8080, State: "LISTEN"},
	}
	bm.SaveIfNeeded(entries)
	bm.SaveIfNeeded(entries) // second call should be a no-op

	b2, _ := scanner.NewBaseline(p)
	if len(b2.Entries()) != 1 {
		t.Errorf("expected 1 baseline entry, got %d", len(b2.Entries()))
	}
}
