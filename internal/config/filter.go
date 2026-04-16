package config

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// FilterRaw holds the raw YAML/JSON filter configuration.
type FilterRaw struct {
	IncludeKeys []string `yaml:"include_keys" json:"include_keys"`
	ExcludeKeys []string `yaml:"exclude_keys" json:"exclude_keys"`
	MinSeverity string   `yaml:"min_severity" json:"min_severity"`
}

var validSeverities = map[string]bool{
	"":         true,
	"low":      true,
	"medium":   true,
	"high":      true,
	"critical": true,
}

// Build validates and converts FilterRaw into a drift.Filter.
func (r FilterRaw) Build() (*drift.Filter, error) {
	if !validSeverities[r.MinSeverity] {
		return nil, fmt.Errorf("invalid min_severity %q: must be one of low, medium, high, critical", r.MinSeverity)
	}
	return &drift.Filter{
		IncludeKeys: r.IncludeKeys,
		ExcludeKeys: r.ExcludeKeys,
		MinSeverity: r.MinSeverity,
	}, nil
}
