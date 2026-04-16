package drift

import "fmt"

// Severity represents the drift severity level.
type Severity string

const (
	SeverityNone     Severity = "none"
	SeverityLow      Severity = "low"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Summary aggregates drift detection statistics.
type Summary struct {
	Total    int      `json:"total"`
	Drifted  int      `json:"drifted"`
	Clean    int      `json:"clean"`
	Severity Severity `json:"severity"`
}

// Summarize computes a Summary from a slice of Results.
func Summarize(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if r.Drift {
			s.Drifted++
		} else {
			s.Clean++
		}
	}
	s.Severity = severityFor(s.Drifted, s.Total)
	return s
}

// String returns a human-readable summary line.
func (s Summary) String() string {
	return fmt.Sprintf("%d/%d keys drifted (severity: %s)", s.Drifted, s.Total, s.Severity)
}

func severityFor(drifted, total int) Severity {
	if total == 0 || drifted == 0 {
		return SeverityNone
	}
	ratio := float64(drifted) / float64(total)
	switch {
	case ratio >= 0.5:
		return SeverityCritical
	case ratio >= 0.2:
		return SeverityHigh
	default:
		return SeverityLow
	}
}
