package config

import (
	"testing"
)

func TestAuditRaw_Build_Disabled(t *testing.T) {
	r := AuditRaw{Enabled: false}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected disabled")
	}
}

func TestAuditRaw_Build_MissingFilePath(t *testing.T) {
	r := AuditRaw{Enabled: true, FilePath: ""}
	_, err := r.Build()
	if err == nil {
		t.Error("expected error for missing file path")
	}
}

func TestAuditRaw_Build_Valid(t *testing.T) {
	r := AuditRaw{Enabled: true, FilePath: "/var/log/audit.json", MaxEntries: 100}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.FilePath != "/var/log/audit.json" {
		t.Errorf("unexpected path: %s", cfg.FilePath)
	}
	if cfg.MaxEntries != 100 {
		t.Errorf("unexpected max entries: %d", cfg.MaxEntries)
	}
}

func TestAuditRaw_Build_DefaultMaxEntries(t *testing.T) {
	r := AuditRaw{Enabled: true, FilePath: "/tmp/audit.json", MaxEntries: 0}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxEntries != 500 {
		t.Errorf("expected default 500, got %d", cfg.MaxEntries)
	}
}
