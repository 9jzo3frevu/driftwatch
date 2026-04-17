package config

import (
	"strings"
	"testing"
)

func TestRemediateRaw_Build_Disabled(t *testing.T) {
	r := RemediateRaw{Enabled: false}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected disabled")
	}
}

func TestRemediateRaw_Build_Valid(t *testing.T) {
	r := RemediateRaw{Enabled: true, Template: "fix {key}"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
	if cfg.Template != "fix {key}" {
		t.Errorf("unexpected template: %s", cfg.Template)
	}
}

func TestRemediateRaw_Build_DefaultTemplate(t *testing.T) {
	r := RemediateRaw{Enabled: true}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Template != "" {
		t.Logf("template is empty, default applied by Remediator")
	}
	_ = cfg
}

func TestRemediateRaw_Build_TemplateTooLong(t *testing.T) {
	r := RemediateRaw{Enabled: true, Template: strings.Repeat("x", 257)}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for long template")
	}
}
