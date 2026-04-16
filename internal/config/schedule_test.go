package config

import (
	"testing"
	"time"
)

func TestScheduleRaw_Build_Disabled(t *testing.T) {
	r := ScheduleRaw{Enabled: false}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Fatal("expected disabled")
	}
}

func TestScheduleRaw_Build_DefaultInterval(t *testing.T) {
	r := ScheduleRaw{Enabled: true}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 60*time.Second {
		t.Fatalf("expected 60s default, got %s", cfg.Interval)
	}
}

func TestScheduleRaw_Build_ValidInterval(t *testing.T) {
	r := ScheduleRaw{Enabled: true, Interval: "5m"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 5*time.Minute {
		t.Fatalf("expected 5m, got %s", cfg.Interval)
	}
	if !cfg.Enabled {
		t.Fatal("expected enabled")
	}
}

func TestScheduleRaw_Build_InvalidInterval(t *testing.T) {
	r := ScheduleRaw{Enabled: true, Interval: "not-a-duration"}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for invalid interval")
	}
}

func TestScheduleRaw_Build_NegativeInterval(t *testing.T) {
	r := ScheduleRaw{Enabled: true, Interval: "-1m"}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for negative interval")
	}
}
