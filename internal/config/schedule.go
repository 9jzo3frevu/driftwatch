package config

import (
	"fmt"
	"time"
)

// ScheduleRaw holds raw schedule config from YAML/JSON.
type ScheduleRaw struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Interval string `yaml:"interval" json:"interval"`
}

// ScheduleConfig is the parsed, validated schedule configuration.
type ScheduleConfig struct {
	Enabled  bool
	Interval time.Duration
}

const defaultScheduleInterval = 60 * time.Second

// Build validates and converts ScheduleRaw into ScheduleConfig.
func (r ScheduleRaw) Build() (ScheduleConfig, error) {
	if !r.Enabled {
		return ScheduleConfig{Enabled: false, Interval: defaultScheduleInterval}, nil
	}
	if r.Interval == "" {
		return ScheduleConfig{Enabled: true, Interval: defaultScheduleInterval}, nil
	}
	d, err := time.ParseDuration(r.Interval)
	if err != nil {
		return ScheduleConfig{}, fmt.Errorf("invalid schedule interval %q: %w", r.Interval, err)
	}
	if d <= 0 {
		return ScheduleConfig{}, fmt.Errorf("schedule interval must be positive, got %s", r.Interval)
	}
	return ScheduleConfig{Enabled: true, Interval: d}, nil
}
