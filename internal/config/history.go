package config

import (
	"errors"
	"time"
)

// HistoryRaw holds raw config for drift history persistence.
type HistoryRaw struct {
	Enabled  bool   `yaml:"enabled"`
	FilePath string `yaml:"file_path"`
	MaxAge   string `yaml:"max_age"`
}

// HistoryConfig is the validated history configuration.
type HistoryConfig struct {
	Enabled  bool
	FilePath string
	MaxAge   time.Duration
}

// Build validates and returns a HistoryConfig.
func (r HistoryRaw) Build() (HistoryConfig, error) {
	if !r.Enabled {
		return HistoryConfig{Enabled: false}, nil
	}
	if r.FilePath == "" {
		return HistoryConfig{}, errors.New("history: file_path is required when enabled")
	}
	maxAge := 7 * 24 * time.Hour
	if r.MaxAge != "" {
		d, err := time.ParseDuration(r.MaxAge)
		if err != nil {
			return HistoryConfig{}, fmt.Errorf("history: invalid max_age %q: %w", r.MaxAge, err)
		}
		if d <= 0 {
			return HistoryConfig{}, errors.New("history: max_age must be positive")
		}
		maxAge = d
	}
	return HistoryConfig{
		Enabled:  true,
		FilePath: r.FilePath,
		MaxAge:   maxAge,
	}, nil
}
