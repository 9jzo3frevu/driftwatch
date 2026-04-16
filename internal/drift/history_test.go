package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func sampleResults(drifted bool) []Result {
	return []Result{
		{
			Key:      "app.replicas",
			Drift:    drifted,
			Declared: ptrStr("3"),
			Actual:   ptrStr("2"),
		},
	}
}

func TestHistory_RecordAndEntries(t *testing.T) {
	dir := t.TempDir()
	h := NewHistory(filepath.Join(dir, "history.json"))

	if err := h.Record(sampleResults(true)); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].DriftCount != 1 {
		t.Errorf("expected DriftCount=1, got %d", entries[0].DriftCount)
	}
}

func TestHistory_NoDriftCount(t *testing.T) {
	dir := t.TempDir()
	h := NewHistory(filepath.Join(dir, "history.json"))

	_ = h.Record(sampleResults(false))
	entries := h.Entries()
	if entries[0].DriftCount != 0 {
		t.Errorf("expected DriftCount=0, got %d", entries[0].DriftCount)
	}
}

func TestHistory_PersistsAndLoads(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h1 := NewHistory(path)
	_ = h1.Record(sampleResults(true))
	_ = h1.Record(sampleResults(false))

	h2 := NewHistory(path)
	if err := h2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(h2.Entries()) != 2 {
		t.Errorf("expected 2 entries after reload, got %d", len(h2.Entries()))
	}
}

func TestHistory_Load_MissingFile(t *testing.T) {
	h := NewHistory("/nonexistent/path/history.json")
	if err := h.Load(); err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestHistory_JSONStructure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	h := NewHistory(path)
	_ = h.Record(sampleResults(true))

	data, _ := os.ReadFile(path)
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid JSON on disk: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry in JSON, got %d", len(entries))
	}
}
