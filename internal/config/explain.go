package config

import "fmt"

// ExplainRaw holds raw configuration for the explainer.
type ExplainRaw struct {
	Enabled bool              `yaml:"enabled"`
	Hints   map[string]string `yaml:"hints"`
}

// ExplainConfig is the validated explainer configuration.
type ExplainConfig struct {
	Enabled bool
	Hints   map[string]string
}

// Build validates and returns an ExplainConfig.
func (r ExplainRaw) Build() (ExplainConfig, error) {
	if !r.Enabled {
		return ExplainConfig{}, nil
	}
	for k, v := range r.Hints {
		if k == "" {
			return ExplainConfig{}, fmt.Errorf("explain: hint key must not be empty")
		}
		if v == "" {
			return ExplainConfig{}, fmt.Errorf("explain: hint value for key %q must not be empty", k)
		}
	}
	hints := r.Hints
	if hints == nil {
		hints = map[string]string{}
	}
	return ExplainConfig{Enabled: true, Hints: hints}, nil
}
