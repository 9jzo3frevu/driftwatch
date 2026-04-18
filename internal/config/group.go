package config

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// GroupRaw holds raw grouping config from YAML.
type GroupRaw struct {
	Enabled bool   `yaml:"enabled"`
	By      string `yaml:"by"`
}

// GroupConfig is the validated grouping configuration.
type GroupConfig struct {
	Enabled bool
	By      drift.GroupKey
}

// Build validates and returns a GroupConfig.
func (r GroupRaw) Build() (GroupConfig, error) {
	if !r.Enabled {
		return GroupConfig{}, nil
	}
	switch drift.GroupKey(r.By) {
	case drift.GroupByService, drift.GroupBySeverity, drift.GroupByKey:
		return GroupConfig{Enabled: true, By: drift.GroupKey(r.By)}, nil
	case "":
		return GroupConfig{Enabled: true, By: drift.GroupByService}, nil
	default:
		return GroupConfig{}, fmt.Errorf("invalid group by value: %q (want service|severity|key)", r.By)
	}
}
