package drift

import "fmt"

// Score represents a numeric drift severity score for a set of results.
type Score struct {
	Total    int
	Critical int
	High     int
	Medium   int
	Low      int
}

// Value returns a weighted composite score.
func (s Score) Value() int {
	return s.Critical*100 + s.High*10 + s.Medium*3 + s.Low*1
}

// String returns a human-readable score summary.
func (s Score) String() string {
	return fmt.Sprintf("score=%d (critical=%d high=%d medium=%d low=%d)",
		s.Value(), s.Critical, s.High, s.Medium, s.Low)
}

// ScoreResults computes a Score from a slice of DriftResult.
func ScoreResults(results []DriftResult) Score {
	s := Score{Total: len(results)}
	for _, r := range results {
		switch r.Severity {
		case "critical":
			s.Critical++
		case "high":
			s.High++
		case "medium":
			s.Medium++
		case "low":
			s.Low++
		}
	}
	return s
}
