package config

import "fmt"

// RemediateRaw holds raw config for the remediation hint feature.
type RemediateRaw struct {
	Enabled  bool   `yaml:"enabled"`
	Template string `yaml:"template"`
}

// RemediateConfig is the validated remediation configuration.
type RemediateConfig struct {
	Enabled  bool
	Template string
}

// Build validates and returns a RemediateConfig.
func (r RemediateRaw) Build() (RemediateConfig, error) {
	if !r.Enabled {
		return RemediateConfig{}, nil
	}
	if len(r.Template) > 256 {
		return RemediateConfig{}, fmt.Errorf("remediate: template exceeds maximum length of 256 characters")
	}
	return RemediateConfig{
		Enabled:  true,
		Template: r.Template,
	}, nil
}
