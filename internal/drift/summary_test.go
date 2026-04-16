package drift

import (
	"strings"
	"testing"
)

func makeResults(total, drifted int) []Result {
	results := make([]Result, total)
	for i := 0; i < total; i++ {
		results[i] = Result{Key: "key", Drift: i < drifted}
	}
	return results
}

func TestSummarize_NoDrift(t *testing.T) {
	s := Summarize(makeResults(5, 0))
	if s.Drifted != 0 || s.Clean != 5 || s.Severity != SeverityNone {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestSummarize_LowSeverity(t *testing.T) {
	s := Summarize(makeResults(10, 1))
	if s.Severity != SeverityLow {
		t.Errorf("expected low, got %s", s.Severity)
	}
}

func TestSummarize_HighSeverity(t *testing.T) {
	s := Summarize(makeResults(10, 3))
	if s.Severity != SeverityHigh {
		t.Errorf("expected high, got %s", s.Severity)
	}
}

func TestSummarize_CriticalSeverity(t *testing.T) {
	s := Summarize(makeResults(10, 6))
	if s.Severity != SeverityCritical {
		t.Errorf("expected critical, got %s", s.Severity)
	}
}

func TestSummarize_EmptyResults(t *testing.T) {
	s := Summarize([]Result{})
	if s.Total != 0 || s.Severity != SeverityNone {
		t.Errorf("unexpected summary for empty: %+v", s)
	}
}

func TestSummary_String(t *testing.T) {
	s := Summarize(makeResults(4, 2))
	str := s.String()
	if !strings.Contains(str, "2/4") {
		t.Errorf("expected '2/4' in string, got: %s", str)
	}
	if !strings.Contains(str, "critical") {
		t.Errorf("expected 'critical' in string, got: %s", str)
	}
}
