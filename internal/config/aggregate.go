package config

import "fmt"

// AggregateRaw holds the raw YAML configuration for result aggregation.
type AggregateRaw struct {
	Enabled        *bool  `yaml:"enabled"`
	GroupBy        string `yaml:"group_by"`
	TopN           int    `yaml:"top_n"`
	IncludeSummary *bool  `yaml:"include_summary"`
}

// AggregateConfig is the validated, ready-to-use aggregation configuration.
type AggregateConfig struct {
	Enabled        bool
	GroupBy        string
	TopN           int
	IncludeSummary bool
}

var validAggregateGroupBy = map[string]struct{}{
	"service":  {},
	"severity": {},
	"prefix":   {},
}

// Build validates and converts AggregateRaw into AggregateConfig.
func (r AggregateRaw) Build() (AggregateConfig, error) {
	if r.Enabled != nil && !*r.Enabled {
		return AggregateConfig{Enabled: false}, nil
	}

	groupBy := r.GroupBy
	if groupBy == "" {
		groupBy = "service"
	}

	if _, ok := validAggregateGroupBy[groupBy]; !ok {
		return AggregateConfig{}, fmt.Errorf(
			"aggregate: invalid group_by %q: must be one of service, severity, prefix", groupBy,
		)
	}

	if r.TopN < 0 {
		return AggregateConfig{}, fmt.Errorf("aggregate: top_n must be >= 0, got %d", r.TopN)
	}

	includeSummary := true
	if r.IncludeSummary != nil {
		includeSummary = *r.IncludeSummary
	}

	return AggregateConfig{
		Enabled:        true,
		GroupBy:        groupBy,
		TopN:           r.TopN,
		IncludeSummary: includeSummary,
	}, nil
}
