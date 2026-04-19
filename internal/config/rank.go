package config

import "fmt"

// RankRaw holds raw config for result ranking.
type RankRaw struct {
	Enabled    bool   `yaml:"enabled"`
	By         string `yaml:"by"`
	Descending bool   `yaml:"descending"`
}

// RankConfig is the validated rank configuration.
type RankConfig struct {
	Enabled    bool
	By         string
	Descending bool
}

var validRankBy = map[string]bool{
	"key":      true,
	"severity": true,
	"service":  true,
}

func (r RankRaw) Build() (RankConfig, error) {
	if !r.Enabled {
		return RankConfig{}, nil
	}
	by := r.By
	if by == "" {
		by = "severity"
	}
	if !validRankBy[by] {
		return RankConfig{}, fmt.Errorf("rank.by must be one of: key, severity, service; got %q", by)
	}
	return RankConfig{
		Enabled:    true,
		By:         by,
		Descending: r.Descending,
	}, nil
}
