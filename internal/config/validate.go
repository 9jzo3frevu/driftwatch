package config

import (
	"fmt"
	"regexp"

	"github.com/driftwatch/internal/drift"
)

// ValidateRuleRaw is the YAML-serialisable form of a validation rule.
type ValidateRuleRaw struct {
	Prefix  string `yaml:"prefix"`
	Pattern string `yaml:"pattern"`
	Message string `yaml:"message"`
}

// Build converts the raw rule into a drift.ValidateRule.
func (r ValidateRuleRaw) Build() (drift.ValidateRule, error) {
	if r.Prefix == "" {
		return drift.ValidateRule{}, fmt.Errorf("validate rule: prefix is required")
	}
	if r.Pattern == "" {
		return drift.ValidateRule{}, fmt.Errorf("validate rule %q: pattern is required", r.Prefix)
	}
	re, err := regexp.Compile(r.Pattern)
	if err != nil {
		return drift.ValidateRule{}, fmt.Errorf("validate rule %q: invalid pattern: %w", r.Prefix, err)
	}
	return drift.ValidateRule{
		Prefix:  r.Prefix,
		Pattern: re,
		Message: r.Message,
	}, nil
}

// BuildValidateRules converts a slice of raw rules into drift.ValidateRule values.
func BuildValidateRules(raws []ValidateRuleRaw) ([]drift.ValidateRule, error) {
	rules := make([]drift.ValidateRule, 0, len(raws))
	for _, raw := range raws {
		r, err := raw.Build()
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}
