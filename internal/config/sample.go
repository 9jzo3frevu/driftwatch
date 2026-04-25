package config

import (
	"fmt"

	"github.com/org/driftwatch/internal/drift"
)

// SampleRaw holds the raw YAML-decoded sampling configuration.
type SampleRaw struct {
	Enabled bool    `yaml:"enabled"`
	Rate    float64 `yaml:"rate"`
	Seed    int64   `yaml:"seed"`
}

// Build validates and converts SampleRaw into a drift.SampleConfig.
// If the section is disabled, the returned config retains all results (rate=1).
func (r *SampleRaw) Build() (drift.SampleConfig, error) {
	if r == nil || !r.Enabled {
		return drift.DefaultSampleConfig(), nil
	}
	if r.Rate <= 0 || r.Rate > 1.0 {
		return drift.SampleConfig{}, fmt.Errorf("config: sample rate must be in (0, 1], got %f", r.Rate)
	}
	return drift.SampleConfig{
		Rate: r.Rate,
		Seed: r.Seed,
	}, nil
}
