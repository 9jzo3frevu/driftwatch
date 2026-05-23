package config

import (
	"strings"
	"testing"
)

func boolPtrCP(b bool) *bool { return &b }

func TestCheckpointRaw_Build_Disabled(t *testing.T) {
	r := &CheckpointRaw{Enabled: boolPtrCP(false), Dir: "/tmp"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected disabled")
	}
}

func TestCheckpointRaw_Build_Nil(t *testing.T) {
	var r *CheckpointRaw
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("nil raw should produce disabled config")
	}
}

func TestCheckpointRaw_Build_MissingDir(t *testing.T) {
	r := &CheckpointRaw{Enabled: boolPtrCP(true), Dir: ""}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for missing dir")
	}
	if !strings.Contains(err.Error(), "dir is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCheckpointRaw_Build_Valid(t *testing.T) {
	r := &CheckpointRaw{Enabled: boolPtrCP(true), Dir: "/var/driftwatch/checkpoints"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
	if cfg.Dir != "/var/driftwatch/checkpoints" {
		t.Errorf("unexpected dir: %s", cfg.Dir)
	}
}

func TestCheckpointRaw_Build_DirTooLong(t *testing.T) {
	r := &CheckpointRaw{Enabled: boolPtrCP(true), Dir: strings.Repeat("a", 257)}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for dir path too long")
	}
	if !strings.Contains(err.Error(), "too long") {
		t.Errorf("unexpected error message: %v", err)
	}
}
