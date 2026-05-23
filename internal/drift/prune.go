package drift

import "time"

// PruneConfig controls how old or excess results are pruned from a result set.
type PruneConfig struct {
	// MaxAge removes results older than this duration (zero disables age pruning).
	MaxAge time.Duration
	// MaxResults caps the total number of results kept (zero disables cap).
	MaxResults int
	// OnlyDrifted removes non-drifted results when true.
	OnlyDrifted bool
}

// DefaultPruneConfig returns a sensible default PruneConfig.
func DefaultPruneConfig() PruneConfig {
	return PruneConfig{
		MaxAge:      24 * time.Hour,
		MaxResults:  500,
		OnlyDrifted: false,
	}
}

// Pruner removes stale or excess DriftResults based on configured rules.
type Pruner struct {
	cfg PruneConfig
	now func() time.Time
}

// NewPruner creates a Pruner with the given config.
func NewPruner(cfg PruneConfig) *Pruner {
	return &Pruner{cfg: cfg, now: time.Now}
}

// Prune applies all configured pruning rules and returns the filtered slice.
func (p *Pruner) Prune(results []DriftResult) []DriftResult {
	out := make([]DriftResult, 0, len(results))
	for _, r := range results {
		if p.cfg.OnlyDrifted && !r.Drifted {
			continue
		}
		if p.cfg.MaxAge > 0 && !r.DetectedAt.IsZero() {
			if p.now().Sub(r.DetectedAt) > p.cfg.MaxAge {
				continue
			}
		}
		out = append(out, r)
	}
	if p.cfg.MaxResults > 0 && len(out) > p.cfg.MaxResults {
		out = out[:p.cfg.MaxResults]
	}
	return out
}
