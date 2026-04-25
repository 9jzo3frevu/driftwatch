package drift

import (
	"testing"
)

func sampleResults() []DriftResult {
	return []DriftResult{
		{Key: "a", IsDrifted: true},
		{Key: "b", IsDrifted: true},
		{Key: "c", IsDrifted: true},
		{Key: "d", IsDrifted: true},
		{Key: "e", IsDrifted: false},
	}
}

func TestSampler_RateOne_KeepsAll(t *testing.T) {
	s := NewSampler(SampleConfig{Rate: 1.0, Seed: 1})
	out := s.Sample(sampleResults())
	if len(out) != 5 {
		t.Fatalf("expected 5 results, got %d", len(out))
	}
}

func TestSampler_RateZero_KeepsNonDrifted(t *testing.T) {
	s := NewSampler(SampleConfig{Rate: 0.0, Seed: 1})
	out := s.Sample(sampleResults())
	// All IsDrifted==true entries should be dropped; IsDrifted==false must stay.
	for _, r := range out {
		if r.IsDrifted {
			t.Errorf("expected no drifted results at rate 0, got key=%s", r.Key)
		}
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 non-drifted result, got %d", len(out))
	}
}

func TestSampler_Deterministic_SameSeed(t *testing.T) {
	cfg := SampleConfig{Rate: 0.5, Seed: 42}
	s1 := NewSampler(cfg)
	s2 := NewSampler(cfg)
	out1 := s1.Sample(sampleResults())
	out2 := s2.Sample(sampleResults())
	if len(out1) != len(out2) {
		t.Fatalf("same seed produced different lengths: %d vs %d", len(out1), len(out2))
	}
	for i := range out1 {
		if out1[i].Key != out2[i].Key {
			t.Errorf("result[%d] mismatch: %s vs %s", i, out1[i].Key, out2[i].Key)
		}
	}
}

func TestSampler_EmptyInput(t *testing.T) {
	s := NewSampler(SampleConfig{Rate: 0.5, Seed: 1})
	out := s.Sample(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}

func TestDefaultSampleConfig_RateIsOne(t *testing.T) {
	cfg := DefaultSampleConfig()
	if cfg.Rate != 1.0 {
		t.Errorf("expected default rate 1.0, got %f", cfg.Rate)
	}
}
