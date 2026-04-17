package config

import "fmt"

// SnapshotRaw holds raw snapshot configuration from the config file.
type SnapshotRaw struct {
	Enabled bool   `yaml:"enabled"`
	Dir     string `yaml:"dir"`
}

// SnapshotConfig is the validated snapshot configuration.
type SnapshotConfig struct {
	Enabled bool
	Dir     string
}

// Build validates and returns a SnapshotConfig.
func (r SnapshotRaw) Build() (SnapshotConfig, error) {
	if !r.Enabled {
		return SnapshotConfig{Enabled: false}, nil
	}
	if r.Dir == "" {
		return SnapshotConfig{}, fmt.Errorf("snapshot: dir must be set when enabled")
	}
	return SnapshotConfig{
		Enabled: true,
		Dir:     r.Dir,
	}, nil
}
