package drift

import "time"

// CapConfig controls result capping behaviour.
type CapConfig struct {
	Enabled    bool
	MaxResults int
	MaxAge     time.Duration
}

// DefaultCapConfig returns sensible defaults.
func DefaultCapConfig() CapConfig {
	return CapConfig{
		Enabled:    true,
		MaxResults: 500,
		MaxAge:     24 * time.Hour,
	}
}

// Cap trims a slice of DriftResult values according to cfg.
// Results are assumed to be ordered most-recent first; entries
// beyond MaxResults or older than MaxAge are dropped.
func Cap(results []DriftResult, cfg CapConfig) []DriftResult {
	if !cfg.Enabled || len(results) == 0 {
		return results
	}

	now := time.Now()
	out := make([]DriftResult, 0, len(results))

	for _, r := range results {
		if cfg.MaxAge > 0 && !r.DetectedAt.IsZero() && now.Sub(r.DetectedAt) > cfg.MaxAge {
			continue
		}
		out = append(out, r)
		if cfg.MaxResults > 0 && len(out) >= cfg.MaxResults {
			break
		}
	}

	return out
}
