package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
)

func strPtr(s string) *string { return &s }

func sampleResults() []drift.DriftResult {
	return []drift.DriftResult{
		{Key: "DB_HOST", ChangeType: drift.Modified, Expected: strPtr("db.prod"), Actual: strPtr("db.staging")},
		{Key: "CACHE_TTL", ChangeType: drift.Removed, Expected: strPtr("300"), Actual: nil},
		{Key: "NEW_FLAG", ChangeType: drift.Added, Expected: nil, Actual: strPtr("true")},
	}
}

func TestFormatter_WriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatText, &buf)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestFormatter_WriteText_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatText, &buf)
	if err := f.Write(sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"DB_HOST", "modified", "db.prod", "db.staging", "CACHE_TTL", "NEW_FLAG"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestFormatter_WriteJSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatJSON, &buf)
	if err := f.Write([]drift.DriftResult{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"drift": []`) {
		t.Errorf("expected empty drift array, got: %s", buf.String())
	}
}

func TestFormatter_WriteJSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatJSON, &buf)
	if err := f.Write(sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{`"DB_HOST"`, `"modified"`, `"db.prod"`, `"db.staging"`} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output:\n%s", want, out)
		}
	}
}
