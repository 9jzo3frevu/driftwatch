package drift

import (
	"strings"
	"unicode"
)

// NormalizeConfig controls how values are normalized before comparison.
type NormalizeConfig struct {
	TrimSpace   bool
	Lowercase   bool
	StripQuotes bool
}

// DefaultNormalizeConfig returns a NormalizeConfig with sensible defaults.
func DefaultNormalizeConfig() NormalizeConfig {
	return NormalizeConfig{
		TrimSpace:   true,
		Lowercase:   false,
		StripQuotes: false,
	}
}

// Normalizer applies normalization rules to drift result values.
type Normalizer struct {
	cfg NormalizeConfig
}

// NewNormalizer creates a Normalizer with the given config.
func NewNormalizer(cfg NormalizeConfig) *Normalizer {
	return &Normalizer{cfg: cfg}
}

// Normalize applies normalization to a copy of each DriftResult.
func (n *Normalizer) Normalize(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	for i, r := range results {
		out[i] = r
		if r.Declared != nil {
			v := n.normalizeValue(*r.Declared)
			out[i].Declared = &v
		}
		if r.Live != nil {
			v := n.normalizeValue(*r.Live)
			out[i].Live = &v
		}
	}
	return out
}

func (n *Normalizer) normalizeValue(v string) string {
	if n.cfg.TrimSpace {
		v = strings.TrimFunc(v, unicode.IsSpace)
	}
	if n.cfg.StripQuotes {
		v = strings.Trim(v, `"'`)
	}
	if n.cfg.Lowercase {
		v = strings.ToLower(v)
	}
	return v
}
