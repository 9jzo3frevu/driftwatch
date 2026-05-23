package config

import "fmt"

var validPartitionBy = map[string]bool{
	"service":    true,
	"severity":   true,
	"key_prefix": true,
}

// PartitionRaw holds raw YAML config for result partitioning.
type PartitionRaw struct {
	Enabled  *bool  `yaml:"enabled"`
	By       string `yaml:"by"`
	MaxSize  int    `yaml:"max_size"`
}

// PartitionConfig is the validated partition configuration.
type PartitionConfig struct {
	Enabled bool
	By      string
	MaxSize int
}

// Build validates and returns a PartitionConfig.
func (r *PartitionRaw) Build() (PartitionConfig, error) {
	if r == nil || (r.Enabled != nil && !*r.Enabled) {
		return PartitionConfig{}, nil
	}

	by := r.By
	if by == "" {
		by = "service"
	}

	if !validPartitionBy[by] {
		return PartitionConfig{}, fmt.Errorf(
			"partition: invalid 'by' value %q: must be one of service, severity, key_prefix", by,
		)
	}

	if r.MaxSize < 0 {
		return PartitionConfig{}, fmt.Errorf("partition: max_size must be >= 0")
	}

	return PartitionConfig{
		Enabled: true,
		By:      by,
		MaxSize: r.MaxSize,
	}, nil
}
