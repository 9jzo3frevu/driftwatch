package drift

import "fmt"

// DedupeConfig controls deduplication behaviour.
type DedupeConfig struct {
	// Window defines how many previous results to compare against.
	Window int
}

// Deduplicator suppresses drift results that were already reported in a
// recent history window, reducing alert noise on persistent drift.
type Deduplicator struct {
	cfg DedupeConfig
}

// NewDeduplicator creates a Deduplicator with the given config.
func NewDeduplicator(cfg DedupeConfig) *Deduplicator {
	if cfg.Window < 1 {
		cfg.Window = 1
	}
	return &Deduplicator{cfg: cfg}
}

// fingerprint returns a stable string key for a DriftResult.
func fingerprint(r DriftResult) string {
	return fmt.Sprintf("%s|%s|%s|%s", r.Service, r.Key, r.Declared, ptrStr(r.Live))
}

// Filter removes results that already appear in any of the provided previous
// entry slices (up to Window entries back).
func (d *Deduplicator) Filter(current []DriftResult, history [][]DriftResult) []DriftResult {
	seen := make(map[string]struct{})

	limit := len(history)
	if limit > d.cfg.Window {
		limit = d.cfg.Window
	}
	for _, entries := range history[len(history)-limit:] {
		for _, r := range entries {
			seen[fingerprint(r)] = struct{}{}
		}
	}

	var out []DriftResult
	for _, r := range current {
		if _, ok := seen[fingerprint(r)]; !ok {
			out = append(out, r)
		}
	}
	return out
}
