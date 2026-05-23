package config

import (
	"fmt"
	"time"

	"github.com/driftwatch/internal/drift"
)

// PruneRaw holds the raw YAML/JSON config for the prune stage.
type PruneRaw struct {
	Enabled     *bool  `yaml:"enabled" json:"enabled"`
	MaxAge      string `yaml:"max_age" json:"max_age"`
	MaxResults  int    `yaml:"max_results" json:"max_results"`
	OnlyDrifted bool   `yaml:"only_drifted" json:"only_drifted"`
}

// Build converts PruneRaw into a drift.PruneConfig.
// Returns nil, nil when disabled.
func (r *PruneRaw) Build() (*drift.PruneConfig, error) {
	if r == nil {
		return nil, nil
	}
	if r.Enabled != nil && !*r.Enabled {
		return nil, nil
	}

	cfg := drift.DefaultPruneConfig()
	cfg.OnlyDrifted = r.OnlyDrifted

	if r.MaxAge != "" {
		d, err := time.ParseDuration(r.MaxAge)
		if err != nil {
			return nil, fmt.Errorf("prune: invalid max_age %q: %w", r.MaxAge, err)
		}
		if d < 0 {
			return nil, fmt.Errorf("prune: max_age must be non-negative")
		}
		cfg.MaxAge = d
	}

	if r.MaxResults > 0 {
		cfg.MaxResults = r.MaxResults
	}

	return &cfg, nil
}
