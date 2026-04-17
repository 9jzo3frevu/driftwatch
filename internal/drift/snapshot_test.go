package drift

import (
	"os"
	"testing"
	"time"
)

func TestNewSnapshot_CopiesValues(t *testing.T) {
	original := map[string]string{"a": "1", "b": "2"}
	s := NewSnapshot("svc", original)
	original["a"] = "mutated"
	if s.Values["a"] != "1" {
		t.Errorf("expected original value, got %s", s.Values["a"])
	}
	if s.ServiceID != "svc" {
		t.Errorf("unexpected service id: %s", s.ServiceID)
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := Snapshot{
		ServiceID:  "api",
		CapturedAt: time.Now().UTC().Truncate(time.Second),
		Values:     map[string]string{"port": "8080", "env": "prod"},
	}
	if err := SaveSnapshot(dir, s); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadSnapshot(dir, "api")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.ServiceID != s.ServiceID {
		t.Errorf("service id mismatch: %s", loaded.ServiceID)
	}
	if loaded.Values["port"] != "8080" {
		t.Errorf("port mismatch: %s", loaded.Values["port"])
	}
}

func TestLoadSnapshot_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadSnapshot(dir, "missing")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/bad_snapshot.json"
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadSnapshot(dir, "bad")
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestSaveSnapshot_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/snapshots"
	s := NewSnapshot("svc", map[string]string{"x": "1"})
	if err := SaveSnapshot(dir, s); err != nil {
		t.Fatalf("expected dir creation, got: %v", err)
	}
}
