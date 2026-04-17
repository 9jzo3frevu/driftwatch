package config

import "testing"

func intPtr(v int) *int { return &v }

func TestScoreRaw_Build_Defaults(t *testing.T) {
	r := ScoreRaw{}
	c := r.Build()
	if c.WarnThreshold != defaultWarnThreshold {
		t.Errorf("expected WarnThreshold=%d, got %d", defaultWarnThreshold, c.WarnThreshold)
	}
	if c.FailThreshold != defaultFailThreshold {
		t.Errorf("expected FailThreshold=%d, got %d", defaultFailThreshold, c.FailThreshold)
	}
}

func TestScoreRaw_Build_CustomValues(t *testing.T) {
	r := ScoreRaw{
		WarnThreshold: intPtr(20),
		FailThreshold: intPtr(200),
	}
	c := r.Build()
	if c.WarnThreshold != 20 {
		t.Errorf("expected WarnThreshold=20, got %d", c.WarnThreshold)
	}
	if c.FailThreshold != 200 {
		t.Errorf("expected FailThreshold=200, got %d", c.FailThreshold)
	}
}

func TestScoreRaw_Build_WarnOnly(t *testing.T) {
	r := ScoreRaw{WarnThreshold: intPtr(5)}
	c := r.Build()
	if c.WarnThreshold != 5 {
		t.Errorf("expected WarnThreshold=5, got %d", c.WarnThreshold)
	}
	if c.FailThreshold != defaultFailThreshold {
		t.Errorf("expected FailThreshold=%d, got %d", defaultFailThreshold, c.FailThreshold)
	}
}
