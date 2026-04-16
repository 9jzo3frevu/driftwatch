package config

import (
	"fmt"
	"time"

	"github.com/driftwatch/internal/drift"
)

// SuppressionRaw holds raw config for a single suppression rule.
type SuppressionRaw struct {
	Key     string `yaml:"key"`
	Expires string `yaml:"expires"`
}

// Build converts a SuppressionRaw into a drift.SuppressionRule.
func (r SuppressionRaw) Build() (drift.SuppressionRule, error) {
	if r.Key == "" {
		return drift.SuppressionRule{}, fmt.Errorf("suppression rule missing key")
	}
	if r.Expires == "" {
		return drift.SuppressionRule{}, fmt.Errorf("suppression rule for %q missing expires", r.Key)
	}
	d, err := time.ParseDuration(r.Expires)
	if err != nil {
		return drift.SuppressionRule{}, fmt.Errorf("invalid expires duration %q: %w", r.Expires, err)
	}
	return drift.SuppressionRule{
		Key:       r.Key,
		ExpiresAt: time.Now().Add(d),
	}, nil
}

// BuildSuppressions converts a slice of raw rules.
func BuildSuppressions(raw []SuppressionRaw) ([]drift.SuppressionRule, error) {
	rules := make([]drift.SuppressionRule, 0, len(raw))
	for _, r := range raw {
		rule, err := r.Build()
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
