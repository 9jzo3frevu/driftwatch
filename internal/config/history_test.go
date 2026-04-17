package config

import (
	"testing"
	"time"
)

func TestHistoryRaw_Build_Disabled(t *testing.T) {
	r := HistoryRaw{Enabled: false}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected disabled")
	}
}

func TestHistoryRaw_Build_MissingFilePath(t *testing.T) {
	r := HistoryRaw{Enabled: true}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for missing file_path")
	}
}

func TestHistoryRaw_Build_DefaultMaxAge(t *testing.T) {
	r := HistoryRaw{Enabled: true, FilePath: "/tmp/history.json"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxAge != 7*24*time.Hour {
		t.Errorf("expected 168h, got %v", cfg.MaxAge)
	}
}

func TestHistoryRaw_Build_ValidMaxAge(t *testing.T) {
	r := HistoryRaw{Enabled: true, FilePath: "/tmp/history.json", MaxAge: "24h"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxAge != 24*time.Hour {
		t.Errorf("expected 24h, got %v", cfg.MaxAge)
	}
}

func TestHistoryRaw_Build_InvalidMaxAge(t *testing.T) {
	r := HistoryRaw{Enabled: true, FilePath: "/tmp/history.json", MaxAge: "notaduration"}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for invalid max_age")
	}
}

func TestHistoryRaw_Build_NegativeMaxAge(t *testing.T) {
	r := HistoryRaw{Enabled: true, FilePath: "/tmp/history.json", MaxAge: "-1h"}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for negative max_age")
	}
}
