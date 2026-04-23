package drift

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func exportResults() []DriftResult {
	decl := "v1"
	live := "v2"
	return []DriftResult{
		{
			Key:      "app.version",
			Service:  "api",
			Declared: &decl,
			Live:     &live,
			Severity: SeverityHigh,
		},
		{
			Key:      "db.pool",
			Service:  "worker",
			Declared: nil,
			Live:     &live,
			Severity: SeverityLow,
		},
	}
}

func TestExporter_CSV_Headers(t *testing.T) {
	e := NewExporter(ExportConfig{Format: ExportCSV})
	var buf bytes.Buffer
	if err := e.Export(&buf, exportResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "key,service") {
		t.Errorf("unexpected header: %s", lines[0])
	}
}

func TestExporter_CSV_Values(t *testing.T) {
	e := NewExporter(ExportConfig{Format: ExportCSV})
	var buf bytes.Buffer
	_ = e.Export(&buf, exportResults())
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.Contains(lines[1], "app.version") {
		t.Errorf("expected key in row: %s", lines[1])
	}
	if !strings.Contains(lines[1], "high") {
		t.Errorf("expected severity in row: %s", lines[1])
	}
}

func TestExporter_CSV_Timestamp(t *testing.T) {
	e := NewExporter(ExportConfig{Format: ExportCSV, Timestamp: true})
	var buf bytes.Buffer
	_ = e.Export(&buf, exportResults())
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.Contains(lines[0], "exported_at") {
		t.Errorf("expected exported_at column in header: %s", lines[0])
	}
}

func TestExporter_JSON_Structure(t *testing.T) {
	e := NewExporter(ExportConfig{Format: ExportJSON})
	var buf bytes.Buffer
	if err := e.Export(&buf, exportResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0]["key"] != "app.version" {
		t.Errorf("unexpected key: %v", rows[0]["key"])
	}
}

func TestExporter_JSON_NoTimestamp(t *testing.T) {
	e := NewExporter(ExportConfig{Format: ExportJSON, Timestamp: false})
	var buf bytes.Buffer
	_ = e.Export(&buf, exportResults())
	var rows []map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &rows)
	if _, ok := rows[0]["exported_at"]; ok {
		t.Error("expected no exported_at field when Timestamp is false")
	}
}

func TestExporter_DefaultFormat(t *testing.T) {
	e := NewExporter(ExportConfig{})
	if e.cfg.Format != ExportCSV {
		t.Errorf("expected default format CSV, got %s", e.cfg.Format)
	}
}

func TestExporter_UnsupportedFormat(t *testing.T) {
	e := NewExporter(ExportConfig{Format: "xml"})
	var buf bytes.Buffer
	err := e.Export(&buf, exportResults())
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
