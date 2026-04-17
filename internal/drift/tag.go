package drift

import "strings"

// Tag represents a label attached to a drift result for grouping or routing.
type Tag struct {
	Key   string
	Value string
}

// Tagger assigns tags to drift results based on key prefix rules.
type Tagger struct {
	rules []TagRule
}

// TagRule maps a key prefix to a set of tags.
type TagRule struct {
	Prefix string
	Tags   []Tag
}

// NewTagger creates a Tagger from the provided rules.
func NewTagger(rules []TagRule) *Tagger {
	return &Tagger{rules: rules}
}

// Apply returns a copy of each DriftResult with matching tags attached.
func (t *Tagger) Apply(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	for i, r := range results {
		r.Tags = t.tagsFor(r.Key)
		out[i] = r
	}
	return out
}

func (t *Tagger) tagsFor(key string) []Tag {
	var matched []Tag
	for _, rule := range t.rules {
		if strings.HasPrefix(key, rule.Prefix) {
			matched = append(matched, rule.Tags...)
		}
	}
	return matched
}
