package config

import (
	"errors"
	"fmt"
)

// CheckpointRaw holds the raw YAML configuration for the checkpoint store.
type CheckpointRaw struct {
	Enabled *bool  `yaml:"enabled"`
	Dir     string `yaml:"dir"`
}

// CheckpointConfig is the validated configuration for the checkpoint store.
type CheckpointConfig struct {
	Enabled bool
	Dir     string
}

// Build validates and returns a CheckpointConfig.
func (r *CheckpointRaw) Build() (CheckpointConfig, error) {
	if r == nil || (r.Enabled != nil && !*r.Enabled) {
		return CheckpointConfig{Enabled: false}, nil
	}
	if r.Dir == "" {
		return CheckpointConfig{}, errors.New("checkpoint: dir is required when enabled")
	}
	if len(r.Dir) > 256 {
		return CheckpointConfig{}, fmt.Errorf("checkpoint: dir path too long (%d chars)", len(r.Dir))
	}
	return CheckpointConfig{
		Enabled: true,
		Dir:     r.Dir,
	}, nil
}
