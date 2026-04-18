package config

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// DedupeRaw holds the raw YAML representation of deduplication settings.
type DedupeRaw struct {
	Enabled *bool `yaml:"enabled"`
	Window  int   `yaml:"window"`
}

// DedupeConfig is the validated, ready-to-use deduplication configuration.
type DedupeConfig struct {
	Enabled bool
	drift.DedupeConfig
}

// Build validates and converts DedupeRaw into a DedupeConfig.
func (r DedupeRaw) Build() (DedupeConfig, error) {
	if r.Enabled != nil && !*r.Enabled {
		return DedupeConfig{Enabled: false}, nil
	}

	win := r.Window
	if win == 0 {
		win = 3 // default look-back window
	}
	if win < 1 || win > 100 {
		return DedupeConfig{}, fmt.Errorf("dedupe window must be between 1 and 100, got %d", win)
	}

	return DedupeConfig{
		Enabled:      true,
		DedupeConfig: drift.DedupeConfig{Window: win},
	}, nil
}
