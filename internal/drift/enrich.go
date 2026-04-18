package drift

import "strings"

// EnrichConfig holds enrichment rules mapping key prefixes to metadata.
type EnrichConfig struct {
	Rules []EnrichRule
}

// EnrichRule associates a key prefix with a set of metadata labels.
type EnrichRule struct {
	Prefix   string
	Metadata map[string]string
}

// Enricher attaches metadata to drift results based on key prefix rules.
type Enricher struct {
	rules []EnrichRule
}

// NewEnricher creates an Enricher from the given config.
func NewEnricher(cfg EnrichConfig) *Enricher {
	return &Enricher{rules: cfg.Rules}
}

// Enrich returns a copy of results with metadata attached where rules match.
func (e *Enricher) Enrich(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	for i, r := range results {
		out[i] = r
		for _, rule := range e.rules {
			if strings.HasPrefix(r.Key, rule.Prefix) {
				if out[i].Metadata == nil {
					out[i].Metadata = make(map[string]string)
				}
				for k, v := range rule.Metadata {
					out[i].Metadata[k] = v
				}
			}
		}
	}
	return out
}
