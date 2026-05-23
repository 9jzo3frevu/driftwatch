package drift

import "time"

// PartitionConfig controls how results are partitioned.
type PartitionConfig struct {
	By      string // "service", "severity", "key_prefix"
	MaxSize int    // max results per partition (0 = unlimited)
}

// DefaultPartitionConfig returns sensible defaults.
func DefaultPartitionConfig() PartitionConfig {
	return PartitionConfig{
		By:      "service",
		MaxSize: 0,
	}
}

// Partition holds a named slice of DriftResults.
type Partition struct {
	Label     string
	Results   []DriftResult
	CreatedAt time.Time
}

// PartitionResults splits results into named partitions based on cfg.By.
// Ordering within each partition preserves input order.
func PartitionResults(results []DriftResult, cfg PartitionConfig) []Partition {
	if len(results) == 0 {
		return nil
	}

	index := make(map[string]int)
	var partitions []Partition
	now := time.Now()

	for _, r := range results {
		label := partitionLabel(r, cfg.By)
		idx, exists := index[label]
		if !exists {
			partitions = append(partitions, Partition{
				Label:     label,
				CreatedAt: now,
			})
			idx = len(partitions) - 1
			index[label] = idx
		}
		p := &partitions[idx]
		if cfg.MaxSize <= 0 || len(p.Results) < cfg.MaxSize {
			p.Results = append(p.Results, r)
		}
	}

	return partitions
}

func partitionLabel(r DriftResult, by string) string {
	switch by {
	case "severity":
		return string(r.Severity)
	case "key_prefix":
		if len(r.Key) > 0 {
			for i, c := range r.Key {
				if c == '.' || c == '_' || c == '/' {
					return r.Key[:i]
				}
			}
		}
		return r.Key
	default: // "service"
		if r.Service != "" {
			return r.Service
		}
		return "unknown"
	}
}
