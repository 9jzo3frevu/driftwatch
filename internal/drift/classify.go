package drift

import (
	"strings"
)

// ClassifyConfig controls how drift results are classified into categories.
type ClassifyConfig struct {
	// CategoryRules maps category names to key prefixes that trigger them.
	CategoryRules map[string][]string
	// DefaultCategory is used when no rule matches.
	DefaultCategory string
}

// DefaultClassifyConfig returns a sensible default classification config.
func DefaultClassifyConfig() ClassifyConfig {
	return ClassifyConfig{
		CategoryRules: map[string][]string{
			"security": {"auth", "secret", "token", "password", "key", "cert"},
			"network":  {"host", "port", "endpoint", "url", "addr"},
			"storage":  {"db", "database", "bucket", "volume", "disk"},
		},
		DefaultCategory: "general",
	}
}

// Classifier assigns a category to each DriftResult based on its key.
type Classifier struct {
	cfg ClassifyConfig
}

// NewClassifier creates a Classifier with the given config.
func NewClassifier(cfg ClassifyConfig) *Classifier {
	return &Classifier{cfg: cfg}
}

// Classify returns a copy of results with the Category field populated.
func (c *Classifier) Classify(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	for i, r := range results {
		r.Category = c.categoryFor(r.Key)
		out[i] = r
	}
	return out
}

func (c *Classifier) categoryFor(key string) string {
	lower := strings.ToLower(key)
	for category, prefixes := range c.cfg.CategoryRules {
		for _, prefix := range prefixes {
			if strings.Contains(lower, prefix) {
				return category
			}
		}
	}
	return c.cfg.DefaultCategory
}
