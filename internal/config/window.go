package config

import (
	"fmt"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// WindowRaw is the YAML-serialisable form of window configuration.
type WindowRaw struct {
	Enabled    bool   `yaml:"enabled"`
	Size       string `yaml:"size"`
	MaxResults int    `yaml:"max_results"`
}

// Build converts WindowRaw into a drift.WindowConfig.
// It returns (cfg, false, nil) when disabled, (cfg, true, nil) when valid,
// and (zero, false, err) on validation failure.
func (r *WindowRaw) Build() (drift.WindowConfig, bool, error) {
	if r == nil || !r.Enabled {
		return drift.WindowConfig{}, false, nil
	}

	defaults := drift.DefaultWindowConfig()
	cfg := drift.WindowConfig{
		Size:       defaults.Size,
		MaxResults: defaults.MaxResults,
	}

	if r.Size != "" {
		d, err := time.ParseDuration(r.Size)
		if err != nil {
			return drift.WindowConfig{}, false, fmt.Errorf("window: invalid size %q: %w", r.Size, err)
		}
		if d <= 0 {
			return drift.WindowConfig{}, false, fmt.Errorf("window: size must be positive, got %s", r.Size)
		}
		cfg.Size = d
	}

	if r.MaxResults > 0 {
		cfg.MaxResults = r.MaxResults
	}

	return cfg, true, nil
}
