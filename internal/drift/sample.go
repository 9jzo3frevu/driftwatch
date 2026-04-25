package drift

import (
	"math/rand"
	"time"
)

// SampleConfig controls probabilistic sampling of drift results.
type SampleConfig struct {
	// Rate is the fraction of results to keep, in the range (0, 1].
	// A rate of 1.0 means keep everything; 0.1 means keep ~10%.
	Rate float64
	// Seed allows deterministic sampling when non-zero (useful in tests).
	Seed int64
}

// DefaultSampleConfig returns a SampleConfig that keeps all results.
func DefaultSampleConfig() SampleConfig {
	return SampleConfig{Rate: 1.0}
}

// Sampler probabilistically retains drift results based on a configured rate.
type Sampler struct {
	cfg SampleConfig
	rng *rand.Rand
}

// NewSampler constructs a Sampler. If cfg.Seed is zero, a time-based seed is
// used so that each run produces independent samples.
func NewSampler(cfg SampleConfig) *Sampler {
	seed := cfg.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &Sampler{
		cfg: cfg,
		//nolint:gosec // non-cryptographic sampling
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Sample filters results, retaining each entry with probability cfg.Rate.
// Results where IsDrifted is false are always retained regardless of rate.
func (s *Sampler) Sample(results []DriftResult) []DriftResult {
	if s.cfg.Rate >= 1.0 {
		return results
	}
	out := make([]DriftResult, 0, len(results))
	for _, r := range results {
		if !r.IsDrifted || s.rng.Float64() < s.cfg.Rate {
			out = append(out, r)
		}
	}
	return out
}
