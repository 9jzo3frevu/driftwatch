package drift

import "strings"

// TruncateConfig holds options for truncating drift result fields.
type TruncateConfig struct {
	MaxValueLen  int
	MaxMessageLen int
}

// DefaultTruncateConfig returns sensible defaults.
func DefaultTruncateConfig() TruncateConfig {
	return TruncateConfig{
		MaxValueLen:   128,
		MaxMessageLen: 256,
	}
}

// Truncator trims long string fields in DriftResults.
type Truncator struct {
	cfg TruncateConfig
}

// NewTruncator creates a Truncator with the given config.
func NewTruncator(cfg TruncateConfig) *Truncator {
	if cfg.MaxValueLen <= 0 {
		cfg.MaxValueLen = DefaultTruncateConfig().MaxValueLen
	}
	if cfg.MaxMessageLen <= 0 {
		cfg.MaxMessageLen = DefaultTruncateConfig().MaxMessageLen
	}
	return &Truncator{cfg: cfg}
}

// Apply returns a new slice of DriftResults with long fields truncated.
func (t *Truncator) Apply(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	for i, r := range results {
		r.Declared = truncStr(r.Declared, t.cfg.MaxValueLen)
		r.Live = truncStr(r.Live, t.cfg.MaxValueLen)
		r.Message = truncStr(r.Message, t.cfg.MaxMessageLen)
		out[i] = r
	}
	return out
}

func truncStr(s *string, max int) *string {
	if s == nil {
		return nil
	}
	if len(*s) <= max {
		return s
	}
	v := (*s)[:max] + "…"
	_ = strings.TrimSpace(v) // ensure package used
	return &v
}
