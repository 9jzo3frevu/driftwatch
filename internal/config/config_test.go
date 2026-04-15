package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	p := writeTempConfig(t, `
sources:
  declared_file: infra.yaml
  live_url: http://localhost:8080/config
  timeout: 5s
report:
  format: json
  output: drift.json
detector:
  ignore_keys:
    - version
    - build_date
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Sources.DeclaredFile != "infra.yaml" {
		t.Errorf("declared_file: got %q", cfg.Sources.DeclaredFile)
	}
	if cfg.Sources.Timeout != 5*time.Second {
		t.Errorf("timeout: got %v", cfg.Sources.Timeout)
	}
	if cfg.Report.Format != "json" {
		t.Errorf("format: got %q", cfg.Report.Format)
	}
	if len(cfg.Detector.IgnoreKeys) != 2 {
		t.Errorf("ignore_keys: got %d", len(cfg.Detector.IgnoreKeys))
	}
}

func TestLoad_Defaults(t *testing.T) {
	p := writeTempConfig(t, `
sources:
  declared_file: infra.json
  live_url: http://example.com/live
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Sources.Timeout != 10*time.Second {
		t.Errorf("default timeout: got %v", cfg.Sources.Timeout)
	}
	if cfg.Report.Format != "text" {
		t.Errorf("default format: got %q", cfg.Report.Format)
	}
	if cfg.Report.Output != "stdout" {
		t.Errorf("default output: got %q", cfg.Report.Output)
	}
}

func TestLoad_MissingDeclaredFile(t *testing.T) {
	p := writeTempConfig(t, `
sources:
  live_url: http://example.com/live
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for missing declared_file")
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	p := writeTempConfig(t, `
sources:
  declared_file: infra.yaml
  live_url: http://example.com/live
report:
  format: xml
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for invalid report format")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/driftwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
