package drift

import "strings"

// Annotation holds metadata attached to a drift result.
type Annotation struct {
	Key   string
	Value string
}

// Annotator attaches annotations to drift results based on key prefix rules.
type Annotator struct {
	rules []AnnotationRule
}

// AnnotationRule maps a key prefix to a set of annotations.
type AnnotationRule struct {
	Prefix      string
	Annotations []Annotation
}

// NewAnnotator creates an Annotator with the given rules.
func NewAnnotator(rules []AnnotationRule) *Annotator {
	return &Annotator{rules: rules}
}

// Annotate returns a copy of results with annotations applied.
func (a *Annotator) Annotate(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	for i, r := range results {
		copy := r
		for _, rule := range a.rules {
			if strings.HasPrefix(r.Key, rule.Prefix) {
				for _, ann := range rule.Annotations {
					copy.Tags = appendUnique(copy.Tags, ann.Key+"="+ann.Value)
				}
			}
		}
		out[i] = copy
	}
	return out
}

func appendUnique(tags []string, tag string) []string {
	for _, t := range tags {
		if t == tag {
			return tags
		}
	}
	return append(tags, tag)
}
