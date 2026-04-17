package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAuditLog_RecordAndEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	al := NewAuditLog(path)

	e := AuditEntry{
		Timestamp:   time.Now().UTC().Truncate(time.Second),
		Service:     "api",
		DriftCount:  3,
		Score:       42,
		TriggeredBy: "scheduler",
	}
	if err := al.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := al.Entries()
	if err != nil {
		t.Fatalf("Entries: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Service != "api" || entries[0].DriftCount != 3 {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestAuditLog_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	al := NewAuditLog(filepath.Join(dir, "audit.jsonl"))

	for i := 0; i < 3; i++ {
		if err := al.Record(AuditEntry{Service: "svc", DriftCount: i}); err != nil {
			t.Fatalf("Record %d: %v", i, err)
		}
	}

	entries, err := al.Entries()
	if err != nil {
		t.Fatalf("Entries: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestAuditLog_Entries_MissingFile(t *testing.T) {
	al := NewAuditLog("/nonexistent/path/audit.jsonl")
	entries, err := al.Entries()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestAuditLog_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "nested", "audit.jsonl")
	al := NewAuditLog(path)
	if err := al.Record(AuditEntry{Service: "x"}); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
