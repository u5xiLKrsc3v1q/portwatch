package monitor

import (
	"testing"

	"github.com/wlynxg/portwatch/internal/scanner"
)

func makeTagEntry(port uint16) scanner.Entry {
	return scanner.Entry{Port: port, Address: "0.0.0.0", Protocol: scanner.TCP}
}

func TestTagFilter_NilTagger_PassesThrough(t *testing.T) {
	f := NewTagFilter(nil, []scanner.Tag{scanner.TagUnknown})
	ev := AlertEvent{Added: []scanner.Entry{makeTagEntry(9000)}}
	out := f.Filter(ev)
	if len(out.Added) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out.Added))
	}
}

func TestTagFilter_SuppressesBaselineEntries(t *testing.T) {
	e := makeTagEntry(443)
	b := &scanner.Baseline{}
	// inject entry manually via exported key method
	b.SetForTest(e)
	tagger := scanner.NewTagger(b, nil)
	f := NewTagFilter(tagger, []scanner.Tag{scanner.TagBaseline})

	ev := AlertEvent{Added: []scanner.Entry{e, makeTagEntry(9999)}}
	out := f.Filter(ev)
	if len(out.Added) != 1 {
		t.Errorf("expected 1 entry after suppression, got %d", len(out.Added))
	}
	if out.Added[0].Port != 9999 {
		t.Errorf("expected port 9999 to remain, got %d", out.Added[0].Port)
	}
}

func TestTagFilter_PreservesRemovedEntries(t *testing.T) {
	tagger := scanner.NewTagger(nil, nil)
	f := NewTagFilter(tagger, []scanner.Tag{scanner.TagNew})
	removed := []scanner.Entry{makeTagEntry(22)}
	ev := AlertEvent{Added: []scanner.Entry{makeTagEntry(8080)}, Removed: removed}
	out := f.Filter(ev)
	if len(out.Removed) != 1 {
		t.Error("removed entries should be preserved")
	}
}

func TestTagFilter_EmptyAdded(t *testing.T) {
	tagger := scanner.NewTagger(nil, nil)
	f := NewTagFilter(tagger, []scanner.Tag{scanner.TagNew})
	ev := AlertEvent{}
	out := f.Filter(ev)
	if len(out.Added) != 0 {
		t.Error("expected empty added")
	}
}
