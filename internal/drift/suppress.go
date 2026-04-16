package drift

import (
	"time"
)

// SuppressionRule defines a rule to suppress drift for a specific key.
type SuppressionRule struct {
	Key       string
	ExpiresAt time.Time
}

// Suppressor holds active suppression rules.
type Suppressor struct {
	rules []SuppressionRule
}

// NewSuppressor creates a new Suppressor.
func NewSuppressor(rules []SuppressionRule) *Suppressor {
	return &Suppressor{rules: rules}
}

// IsSuppressed returns true if the given key has an active suppression rule.
func (s *Suppressor) IsSuppressed(key string, now time.Time) bool {
	for _, r := range s.rules {
		if r.Key == key && now.Before(r.ExpiresAt) {
			return true
		}
	}
	return false
}

// Apply filters out suppressed drift results.
func (s *Suppressor) Apply(results []DriftResult, now time.Time) []DriftResult {
	filtered := make([]DriftResult, 0, len(results))
	for _, r := range results {
		if !s.IsSuppressed(r.Key, now) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
