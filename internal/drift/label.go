package drift

import (
	"strings"
)

// LabelConfig controls how drift results are labelled.
type LabelConfig struct {
	// Rules maps key prefixes to a set of label key=value pairs.
	Rules []LabelRule
}

// LabelRule associates a key prefix with labels to apply.
type LabelRule struct {
	Prefix string
	Labels map[string]string
}

// Labeler attaches metadata labels to drift results based on key prefix rules.
type Labeler struct {
	rules []LabelRule
}

// NewLabeler constructs a Labeler from the provided config.
func NewLabeler(cfg LabelConfig) *Labeler {
	return &Labeler{rules: cfg.Rules}
}

// Apply returns a copy of results with labels merged in according to matching rules.
func (l *Labeler) Apply(results []DriftResult) []DriftResult {
	if len(l.rules) == 0 {
		return results
	}
	out := make([]DriftResult, len(results))
	for i, r := range results {
		copy := r
		for _, rule := range l.rules {
			if strings.HasPrefix(r.Key, rule.Prefix) {
				if copy.Labels == nil {
					copy.Labels = make(map[string]string)
				}
				for k, v := range rule.Labels {
					copy.Labels[k] = v
				}
			}
		}
		out[i] = copy
	}
	return out
}
