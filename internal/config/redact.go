package config

import (
	"github.com/driftwatch/internal/drift"
)

// RedactRaw holds raw config for value redaction.
type RedactRaw struct {
	Enabled  bool     `yaml:"enabled"`
	Patterns []string `yaml:"patterns"`
	Mask     string   `yaml:"mask"`
}

// Build validates and constructs a drift.RedactConfig.
func (r RedactRaw) Build() (drift.RedactConfig, bool, error) {
	if !r.Enabled {
		return drift.RedactConfig{}, false, nil
	}
	cfg := drift.DefaultRedactConfig()
	if len(r.Patterns) > 0 {
		cfg.Patterns = r.Patterns
	}
	if r.Mask != "" {
		cfg.Mask = r.Mask
	}
	return cfg, true, nil
}
