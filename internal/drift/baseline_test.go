package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewBaseline_CopiesValues(t *testing.T) {
	orig := map[string]string{"a": "1", "b": "2"}
	b := NewBaseline(orig)
	orig["a"] = "mutated"
	if b.Values["a"] != "1" {
		t.Errorf("expected baseline to be immutable, got %s", b.Values["a"])
	}
	if b.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	original := NewBaseline(map[string]string{"host": "localhost", "port": "8080"})
	if err := SaveBaseline(path, original); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	if loaded.Values["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %s", loaded.Values["host"])
	}
	if loaded.Values["port"] != "8080" {
		t.Errorf("expected port=8080, got %s", loaded.Values["port"])
	}
	if !loaded.CapturedAt.Equal(original.CapturedAt) {
		t.Errorf("CapturedAt mismatch: %v vs %v", loaded.CapturedAt, original.CapturedAt)
	}
}

func TestLoadBaseline_NotFound(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := LoadBaseline(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
