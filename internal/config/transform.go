package config

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// TransformRuleRaw is the raw YAML representation of a transform rule.
type TransformRuleRaw struct {
	Prefix string `yaml:"prefix"`
	Op     string `yaml:"op"`
	Arg    string `yaml:"arg"`
}

// TransformRaw is the raw YAML block for the transform section.
type TransformRaw struct {
	Enabled bool               `yaml:"enabled"`
	Rules   []TransformRuleRaw `yaml:"rules"`
}

var validTransformOps = map[string]drift.TransformOp{
	"prepend": drift.TransformOpPrefix,
	"append":  drift.TransformOpSuffix,
	"replace": drift.TransformOpReplace,
	"upper":   drift.TransformOpUpper,
	"lower":   drift.TransformOpLower,
}

// Build validates and converts TransformRaw into a drift.TransformConfig.
// Returns nil config and no error when disabled.
func (r TransformRaw) Build() (*drift.TransformConfig, error) {
	if !r.Enabled {
		return nil, nil
	}
	var rules []drift.TransformRule
	for i, raw := range r.Rules {
		if raw.Prefix == "" {
			return nil, fmt.Errorf("transform rule %d: prefix is required", i)
		}
		op, ok := validTransformOps[raw.Op]
		if !ok {
			return nil, fmt.Errorf("transform rule %d: unknown op %q", i, raw.Op)
		}
		rules = append(rules, drift.TransformRule{
			Prefix: raw.Prefix,
			Op:     op,
			Arg:    raw.Arg,
		})
	}
	return &drift.TransformConfig{Rules: rules}, nil
}

// BuildTransformRules is a convenience wrapper for building rules from a slice.
func BuildTransformRules(rules []TransformRuleRaw) ([]drift.TransformRule, error) {
	tmp := TransformRaw{Enabled: true, Rules: rules}
	cfg, err := tmp.Build()
	if err != nil {
		return nil, err
	}
	return cfg.Rules, nil
}
