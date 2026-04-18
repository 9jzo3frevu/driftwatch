package config

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// EnrichRuleRaw holds raw config for a single enrichment rule.
type EnrichRuleRaw struct {
	Prefix   string            `yaml:"prefix"`
	Metadata map[string]string `yaml:"metadata"`
}

// Build validates and converts an EnrichRuleRaw to drift.EnrichRule.
func (r EnrichRuleRaw) Build() (drift.EnrichRule, error) {
	if r.Prefix == "" {
		return drift.EnrichRule{}, fmt.Errorf("enrich rule missing prefix")
	}
	if len(r.Metadata) == 0 {
		return drift.EnrichRule{}, fmt.Errorf("enrich rule for prefix %q has no metadata", r.Prefix)
	}
	return drift.EnrichRule{Prefix: r.Prefix, Metadata: r.Metadata}, nil
}

// BuildEnrichRules converts a slice of raw rules into drift.EnrichConfig.
func BuildEnrichRules(raw []EnrichRuleRaw) (drift.EnrichConfig, error) {
	rules := make([]drift.EnrichRule, 0, len(raw))
	for _, r := range raw {
		rule, err := r.Build()
		if err != nil {
			return drift.EnrichConfig{}, err
		}
		rules = append(rules, rule)
	}
	return drift.EnrichConfig{Rules: rules}, nil
}
