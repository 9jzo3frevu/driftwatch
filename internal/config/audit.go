package config

import "fmt"

// AuditRaw holds raw audit log configuration from YAML.
type AuditRaw struct {
	Enabled    bool   `yaml:"enabled"`
	FilePath   string `yaml:"file_path"`
	MaxEntries int    `yaml:"max_entries"`
}

// AuditConfig is the validated audit log configuration.
type AuditConfig struct {
	Enabled    bool
	FilePath   string
	MaxEntries int
}

const defaultAuditMaxEntries = 500

// Build validates and returns an AuditConfig.
func (r AuditRaw) Build() (AuditConfig, error) {
	if !r.Enabled {
		return AuditConfig{Enabled: false}, nil
	}
	if r.FilePath == "" {
		return AuditConfig{}, fmt.Errorf("audit: file_path is required when enabled")
	}
	max := r.MaxEntries
	if max <= 0 {
		max = defaultAuditMaxEntries
	}
	return AuditConfig{
		Enabled:    true,
		FilePath:   r.FilePath,
		MaxEntries: max,
	}, nil
}
