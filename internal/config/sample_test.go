package config

import (
	"testing"
)

func TestSampleRaw_Build_Disabled(t *testing.T) {
	r := &SampleRaw{Enabled: false, Rate: 0.5}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Rate != 1.0 {
		t.Errorf("expected rate 1.0 when disabled, got %f", cfg.Rate)
	}
}

func TestSampleRaw_Build_NilIsDisabled(t *testing.T) {
	var r *SampleRaw
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Rate != 1.0 {
		t.Errorf("expected rate 1.0 for nil config, got %f", cfg.Rate)
	}
}

func TestSampleRaw_Build_Valid(t *testing.T) {
	r := &SampleRaw{Enabled: true, Rate: 0.25, Seed: 99}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Rate != 0.25 {
		t.Errorf("expected rate 0.25, got %f", cfg.Rate)
	}
	if cfg.Seed != 99 {
		t.Errorf("expected seed 99, got %d", cfg.Seed)
	}
}

func TestSampleRaw_Build_RateZero_Error(t *testing.T) {
	r := &SampleRaw{Enabled: true, Rate: 0.0}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for rate=0, got nil")
	}
}

func TestSampleRaw_Build_RateAboveOne_Error(t *testing.T) {
	r := &SampleRaw{Enabled: true, Rate: 1.5}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for rate>1, got nil")
	}
}

func TestSampleRaw_Build_RateExactlyOne_Valid(t *testing.T) {
	r := &SampleRaw{Enabled: true, Rate: 1.0}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Rate != 1.0 {
		t.Errorf("expected rate 1.0, got %f", cfg.Rate)
	}
}
