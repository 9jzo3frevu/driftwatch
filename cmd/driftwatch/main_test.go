package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestRun_MissingConfig(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"driftwatch", "/nonexistent/path/config.yaml"}

	err := run()
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestRun_InvalidConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := writeTempFile(t, dir, "bad.yaml", ":::invalid yaml:::")

	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"driftwatch", cfgPath}

	err := run()
	if err == nil {
		t.Fatal("expected error for invalid config, got nil")
	}
}

func TestRun_DefaultConfigPath(t *testing.T) {
	// Ensure run() attempts default config path when no args given.
	// It will fail loading the config (file not present in test env),
	// but the error should reference loading config, not a panic.
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"driftwatch"}

	// Change working dir to temp so default driftwatch.yaml is absent.
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir) //nolint:errcheck
	os.Chdir(t.TempDir())   //nolint:errcheck

	err := run()
	if err == nil {
		t.Fatal("expected error when default config absent, got nil")
	}
}
