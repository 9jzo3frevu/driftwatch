package config

import "fmt"

// PipelineStageRaw is the raw YAML representation of a single pipeline stage.
type PipelineStageRaw struct {
	Name    string `yaml:"name"`
	Enabled *bool  `yaml:"enabled"`
}

// PipelineRaw is the raw YAML representation of the pipeline configuration.
type PipelineRaw struct {
	Stages []PipelineStageRaw `yaml:"stages"`
}

// PipelineConfig is the validated pipeline configuration.
type PipelineConfig struct {
	// EnabledStages is the ordered list of stage names to include.
	EnabledStages []string
}

// Build validates and converts PipelineRaw into a PipelineConfig.
func (r *PipelineRaw) Build() (*PipelineConfig, error) {
	if r == nil || len(r.Stages) == 0 {
		return &PipelineConfig{}, nil
	}

	seen := make(map[string]bool)
	var enabled []string

	for i, s := range r.Stages {
		if s.Name == "" {
			return nil, fmt.Errorf("pipeline stage[%d]: name is required", i)
		}
		if seen[s.Name] {
			return nil, fmt.Errorf("pipeline stage[%d]: duplicate stage name %q", i, s.Name)
		}
		seen[s.Name] = true
		if s.Enabled != nil && !*s.Enabled {
			continue
		}
		enabled = append(enabled, s.Name)
	}

	return &PipelineConfig{EnabledStages: enabled}, nil
}
