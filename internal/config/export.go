package config

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// ExportRaw holds the raw YAML fields for export configuration.
type ExportRaw struct {
	Enabled   bool   `yaml:"enabled"`
	Format    string `yaml:"format"`
	Timestamp bool   `yaml:"timestamp"`
}

// Build validates and converts ExportRaw into a drift.ExportConfig.
// Returns nil, nil when export is disabled.
func (r ExportRaw) Build() (*drift.ExportConfig, error) {
	if !r.Enabled {
		return nil, nil
	}

	fmt := drift.ExportFormat(r.Format)
	switch fmt {
	case "", drift.ExportCSV:
		fmt = drift.ExportCSV
	case drift.ExportJSON:
		// valid
	default:
		return nil, fmt.Errorf("export: unsupported format %q; must be \"csv\" or \"json\"", r.Format)
	}

	return &drift.ExportConfig{
		Format:    fmt,
		Timestamp: r.Timestamp,
	}, nil
}
