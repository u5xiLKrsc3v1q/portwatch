package scanner

// Tag represents a label attached to a port entry for classification.
type Tag string

const (
	TagKnown    Tag = "known"
	TagUnknown  Tag = "unknown"
	TagBaseline Tag = "baseline"
	TagNew      Tag = "new"
)

// Tagger assigns tags to entries based on configurable rules.
type Tagger struct {
	baseline  *Baseline
	whitelist *Whitelist
}

// NewTagger creates a Tagger. Both baseline and whitelist may be nil.
func NewTagger(b *Baseline, w *Whitelist) *Tagger {
	return &Tagger{baseline: b, whitelist: w}
}

// Tag returns the set of tags applicable to the given entry.
func (t *Tagger) Tag(e Entry) []Tag {
	var tags []Tag

	if t.baseline != nil && t.baseline.Contains(e) {
		tags = append(tags, TagBaseline)
	} else {
		tags = append(tags, TagNew)
	}

	if t.whitelist != nil && t.whitelist.IsAllowed(e) {
		tags = append(tags, TagKnown)
	} else {
		tags = append(tags, TagUnknown)
	}

	return tags
}

// HasTag reports whether tags contains the given tag.
func HasTag(tags []Tag, target Tag) bool {
	for _, tg := range tags {
		if tg == target {
			return true
		}
	}
	return false
}
