package drift

import (
	"testing"
)

func scoreResults() []DriftResult {
	return []DriftResult{
		{Key: "a", Severity: "critical"},
		{Key: "b", Severity: "high"},
		{Key: "c", Severity: "medium"},
		{Key: "d", Severity: "low"},
		{Key: "e", Severity: "low"},
	}
}

func TestScoreResults_Counts(t *testing.T) {
	s := ScoreResults(scoreResults())
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.Critical != 1 {
		t.Errorf("expected Critical=1, got %d", s.Critical)
	}
	if s.High != 1 {
		t.Errorf("expected High=1, got %d", s.High)
	}
	if s.Medium != 1 {
		t.Errorf("expected Medium=1, got %d", s.Medium)
	}
	if s.Low != 2 {
		t.Errorf("expected Low=2, got %d", s.Low)
	}
}

func TestScoreResults_Value(t *testing.T) {
	s := ScoreResults(scoreResults())
	// 1*100 + 1*10 + 1*3 + 2*1 = 115
	if s.Value() != 115 {
		t.Errorf("expected Value=115, got %d", s.Value())
	}
}

func TestScoreResults_Empty(t *testing.T) {
	s := ScoreResults(nil)
	if s.Value() != 0 {
		t.Errorf("expected Value=0, got %d", s.Value())
	}
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
}

func TestScore_String(t *testing.T) {
	s := Score{Total: 1, Critical: 1}
	got := s.String()
	if got == "" {
		t.Error("expected non-empty string")
	}
}
