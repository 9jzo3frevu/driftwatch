package config

import (
	"fmt"
	"path/filepath"
)

// AuditRaw holds raw config for the audit log feature.
type AuditRaw struct {
	Enabled  bool   `yaml:"enabled"`
	FilePath string `yaml:"file_path"`
}

// AuditConfig is the validated audit log configuration.
type AuditConfig struct {
	Enabled  bool
	FilePath string
}

// Build validates and returns an AuditConfig.
func (r AuditRaw) Build() (AuditConfig, error) {
	if !r.Enabled {
		return AuditConfig{}, nil
	}
	if r.FilePath == "" {
		return AuditConfig{}, fmt.Errorf("audit: file_path is required when enabled")
	}
	if !filepath.IsAbs(r.FilePath) && filepath.Ext(r.FilePath) == "" {
		return AuditConfig{}, fmt.Errorf("audit: file_path must include a filename")
	}
	return AuditConfig{
		Enabled:  true,
		FilePath: r.FilePath,
	}, nil
}
