package config

import (
	"fmt"

	"github.com/org/driftwatch/internal/drift"
)

// TagRuleRaw holds raw config for a single tag rule.
type TagRuleRaw struct {
	Prefix string            `yaml:"prefix"`
	Tags   map[string]string `yaml:"tags"`
}

// Build converts a TagRuleRaw into a drift.TagRule.
func (r TagRuleRaw) Build() (drift.TagRule, error) {
	if r.Prefix == "" {
		return drift.TagRule{}, fmt.Errorf("tag rule missing prefix")
	}
	if len(r.Tags) == 0 {
		return drift.TagRule{}, fmt.Errorf("tag rule for prefix %q has no tags", r.Prefix)
	}
	var tags []drift.Tag
	for k, v := range r.Tags {
		tags = append(tags, drift.Tag{Key: k, Value: v})
	}
	return drift.TagRule{Prefix: r.Prefix, Tags: tags}, nil
}

// BuildTagRules converts a slice of raw rules into drift.TagRules.
func BuildTagRules(raws []TagRuleRaw) ([]drift.TagRule, error) {
	var rules []drift.TagRule
	for _, raw := range raws {
		r, err := raw.Build()
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}
