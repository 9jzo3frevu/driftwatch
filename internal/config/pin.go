package config

import (
	"fmt"
	"time"
)

// PinRaw is the YAML-decoded representation of the pin configuration block.
type PinRaw struct {
	Enabled  *bool  `yaml:"enabled"`
	FilePath string `yaml:"file_path"`
}

// PinConfig holds the validated pin store configuration.
type PinConfig struct {
	Enabled  bool
	FilePath string
}

// Build validates and converts PinRaw into a PinConfig.
func (r *PinRaw) Build() (PinConfig, error) {
	if r == nil || r.Enabled == nil || !*r.Enabled {
		return PinConfig{}, nil
	}
	if r.FilePath == "" {
		return PinConfig{}, fmt.Errorf("pin: file_path is required when enabled")
	}
	return PinConfig{
		Enabled:  true,
		FilePath: r.FilePath,
	}, nil
}

// PinSuppressionRaw represents a single key-level pin defined in YAML.
type PinSuppressionRaw struct {
	Key     string `yaml:"key"`
	Service string `yaml:"service"`
	Expires string `yaml:"expires,omitempty"`
	Reason  string `yaml:"reason,omitempty"`
}

// PinSuppression is the validated form of a PinSuppressionRaw.
type PinSuppression struct {
	Key      string
	Service  string
	Expires  time.Duration
	Reason   string
}

// Build validates and converts a PinSuppressionRaw.
func (r PinSuppressionRaw) Build() (PinSuppression, error) {
	if r.Key == "" {
		return PinSuppression{}, fmt.Errorf("pin suppression: key is required")
	}
	var dur time.Duration
	if r.Expires != "" {
		var err error
		dur, err = time.ParseDuration(r.Expires)
		if err != nil {
			return PinSuppression{}, fmt.Errorf("pin suppression: invalid expires %q: %w", r.Expires, err)
		}
		if dur <= 0 {
			return PinSuppression{}, fmt.Errorf("pin suppression: expires must be positive")
		}
	}
	return PinSuppression{
		Key:     r.Key,
		Service: r.Service,
		Expires: dur,
		Reason:  r.Reason,
	}, nil
}
