package scanner

import (
	"testing"
)

func makeTaggerEntry(port uint16, addr string) Entry {
	return Entry{Port: port, Address: addr, Protocol: TCP}
}

func TestTagger_NilBaselineAndWhitelist(t *testing.T) {
	tagger := NewTagger(nil, nil)
	e := makeTaggerEntry(8080, "0.0.0.0")
	tags := tagger.Tag(e)

	if !HasTag(tags, TagNew) {
		t.Error("expected TagNew when no baseline")
	}
	if !HasTag(tags, TagUnknown) {
		t.Error("expected TagUnknown when no whitelist")
	}
}

func TestTagger_BaselineEntry(t *testing.T) {
	e := makeTaggerEntry(443, "0.0.0.0")
	b := &Baseline{entries: map[string]struct{}{e.Key(): {}}}
	tagger := NewTagger(b, nil)
	tags := tagger.Tag(e)

	if !HasTag(tags, TagBaseline) {
		t.Error("expected TagBaseline for known baseline entry")
	}
	if HasTag(tags, TagNew) {
		t.Error("did not expect TagNew for baseline entry")
	}
}

func TestTagger_WhitelistedEntry(t *testing.T) {
	e := makeTaggerEntry(80, "0.0.0.0")
	w := NewWhitelist([]WhitelistRule{{Port: 80}})
	tagger := NewTagger(nil, w)
	tags := tagger.Tag(e)

	if !HasTag(tags, TagKnown) {
		t.Error("expected TagKnown for whitelisted entry")
	}
	if HasTag(tags, TagUnknown) {
		t.Error("did not expect TagUnknown for whitelisted entry")
	}
}

func TestHasTag_True(t *testing.T) {
	tags := []Tag{TagNew, TagKnown}
	if !HasTag(tags, TagKnown) {
		t.Error("expected HasTag to return true")
	}
}

func TestHasTag_False(t *testing.T) {
	tags := []Tag{TagNew}
	if HasTag(tags, TagBaseline) {
		t.Error("expected HasTag to return false")
	}
}
