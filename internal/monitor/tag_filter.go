package monitor

import (
	"github.com/wlynxg/portwatch/internal/scanner"
)

// TagFilter suppresses alert events whose added entries carry suppressed tags.
type TagFilter struct {
	tagger     *scanner.Tagger
	suppress   []scanner.Tag
}

// NewTagFilter creates a TagFilter that will drop added entries having any of
// the given suppress tags according to the provided Tagger.
func NewTagFilter(tagger *scanner.Tagger, suppress []scanner.Tag) *TagFilter {
	return &TagFilter{tagger: tagger, suppress: suppress}
}

// Filter returns a new AlertEvent with suppressed entries removed from Added.
// Removed entries are always preserved.
func (f *TagFilter) Filter(ev AlertEvent) AlertEvent {
	if f.tagger == nil {
		return ev
	}
	var kept []scanner.Entry
	for _, e := range ev.Added {
		tags := f.tagger.Tag(e)
		if !f.anySuppressed(tags) {
			kept = append(kept, e)
		}
	}
	return AlertEvent{Added: kept, Removed: ev.Removed}
}

func (f *TagFilter) anySuppressed(tags []scanner.Tag) bool {
	for _, s := range f.suppress {
		if scanner.HasTag(tags, s) {
			return true
		}
	}
	return false
}
