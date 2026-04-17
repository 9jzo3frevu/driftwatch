package config

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// DiffRaw holds the raw config for diff behaviour.
type DiffRaw struct {
	Mode string `yaml:"mode"` // "exact" or "subset"
}

// DiffConfig holds the resolved diff configuration.
type DiffConfig struct {
	Mode drift.DiffMode
}

// Build validates and resolves DiffRaw into DiffConfig.
func (r DiffRaw) Build() (DiffConfig, error) {
	switch r.Mode {
	case "", "exact":
		return DiffConfig{Mode: drift.DiffModeExact}, nil
	case "subset":
		return DiffConfig{Mode: drift.DiffModeSubset}, nil
	default:
		return DiffConfig{}, fmt.Errorf("unknown diff mode %q: must be \"exact\" or \"subset\"", r.Mode)
	}
}
