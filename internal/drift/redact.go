package drift

import (
	"strings"
)

// RedactConfig controls which keys should have their values redacted.
type RedactConfig struct {
	Patterns []string
	Mask     string
}

// DefaultRedactConfig returns a RedactConfig with sensible defaults.
func DefaultRedactConfig() RedactConfig {
	return RedactConfig{
		Patterns: []string{"password", "secret", "token", "key", "apikey", "api_key"},
		Mask:     "***REDACTED***",
	}
}

// Redactor masks sensitive values in drift results.
type Redactor struct {
	cfg RedactConfig
}

// NewRedactor creates a Redactor from the given config.
func NewRedactor(cfg RedactConfig) *Redactor {
	if cfg.Mask == "" {
		cfg.Mask = "***REDACTED***"
	}
	return &Redactor{cfg: cfg}
}

// Redact returns a copy of results with sensitive values masked.
func (r *Redactor) Redact(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	for i, res := range results {
		out[i] = res
		if r.isSensitive(res.Key) {
			if out[i].Declared != nil {
				v := r.cfg.Mask
				out[i].Declared = &v
			}
			if out[i].Live != nil {
				v := r.cfg.Mask
				out[i].Live = &v
			}
		}
	}
	return out
}

func (r *Redactor) isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.cfg.Patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
