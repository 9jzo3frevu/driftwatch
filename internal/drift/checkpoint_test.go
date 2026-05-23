package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCheckpointStore_SaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := NewCheckpointStore(dir)

	cp := Checkpoint{
		Service:    "api",
		RunAt:      time.Now().UTC().Truncate(time.Second),
		DriftCount: 3,
	}
	if err := store.Save(cp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load("api")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Service != cp.Service || got.DriftCount != cp.DriftCount || !got.RunAt.Equal(cp.RunAt) {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, cp)
	}
}

func TestCheckpointStore_Load_NotFound(t *testing.T) {
	store := NewCheckpointStore(t.TempDir())
	cp, err := store.Load("nonexistent")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cp.Service != "" {
		t.Errorf("expected zero checkpoint, got %+v", cp)
	}
}

func TestCheckpointStore_Load_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.checkpoint.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	store := NewCheckpointStore(dir)
	_, err := store.Load("bad")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCheckpointStore_Delete(t *testing.T) {
	dir := t.TempDir()
	store := NewCheckpointStore(dir)

	cp := Checkpoint{Service: "svc", RunAt: time.Now(), DriftCount: 1}
	_ = store.Save(cp)

	if err := store.Delete("svc"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, err := store.Load("svc")
	if err != nil || got.Service != "" {
		t.Errorf("expected zero checkpoint after delete, got %+v err %v", got, err)
	}
}

func TestCheckpointStore_Delete_NotFound(t *testing.T) {
	store := NewCheckpointStore(t.TempDir())
	if err := store.Delete("ghost"); err != nil {
		t.Fatalf("Delete of nonexistent should not error: %v", err)
	}
}

func TestCheckpointStore_Save_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "checkpoints")
	store := NewCheckpointStore(dir)
	cp := Checkpoint{Service: "x", RunAt: time.Now(), DriftCount: 0}
	if err := store.Save(cp); err != nil {
		t.Fatalf("Save should create dirs: %v", err)
	}
}
