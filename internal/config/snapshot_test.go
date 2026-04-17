package config

import "testing"

func TestSnapshotRaw_Build_Disabled(t *testing.T) {
	r := SnapshotRaw{Enabled: false}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected disabled")
	}
}

func TestSnapshotRaw_Build_MissingDir(t *testing.T) {
	r := SnapshotRaw{Enabled: true, Dir: ""}
	_, err := r.Build()
	if err == nil {
		t.Error("expected error for missing dir")
	}
}

func TestSnapshotRaw_Build_Valid(t *testing.T) {
	r := SnapshotRaw{Enabled: true, Dir: "/tmp/snapshots"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
	if cfg.Dir != "/tmp/snapshots" {
		t.Errorf("unexpected dir: %s", cfg.Dir)
	}
}
