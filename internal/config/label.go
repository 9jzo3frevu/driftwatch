package config

import (
	"fmt"

	"github.com/your-org/driftwatch/internal/drift"
)

// LabelRuleRaw is the YAML-decoded representation of a single label rule.
type LabelRuleRaw struct {
	Prefix string            `yaml:"prefix"`
	Labels map[string]string `yaml:"labels"`
}

// Build validates and converts the raw rule into a drift.LabelRule.
func (r LabelRuleRaw) Build() (drift.LabelRule, error) {
	if r.Prefix == "" {
		return drift.LabelRule{}, fmt.Errorf("label rule: prefix must not be empty")
	}
	if len(r.Labels) == 0 {
		return drift.LabelRule{}, fmt.Errorf("label rule: labels map must not be empty for prefix %q", r.Prefix)
	}
	return drift.LabelRule{Prefix: r.Prefix, Labels: r.Labels}, nil
}

// BuildLabelRules converts a slice of raw rules into drift.LabelRule values.
func BuildLabelRules(raw []LabelRuleRaw) ([]drift.LabelRule, error) {
	rules := make([]drift.LabelRule, 0, len(raw))
	for _, r := range raw {
		rule, err := r.Build()
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
