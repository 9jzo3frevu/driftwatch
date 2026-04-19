package drift

import (
	"strings"
	"testing"
)

func truncResults() []DriftResult {
	return []DriftResult{
		{
			Key:      "db.password",
			Declared: ptrStr(strings.Repeat("a", 200)),
			Live:     ptrStr(strings.Repeat("b", 200)),
			Message:  strings.Repeat("c", 300),
		},
		{
			Key:      "app.name",
			Declared: ptrStr("short"),
			Live:     ptrStr("short"),
			Message:  "ok",
		},
	}
}

func TestTruncator_LongValues(t *testing.T) {
	tr := NewTruncator(DefaultTruncateConfig())
	out := tr.Apply(truncResults())

	if len(*out[0].Declared) > 129 {
		t.Errorf("declared not truncated: len=%d", len(*out[0].Declared))
	}
	if len(*out[0].Live) > 129 {
		t.Errorf("live not truncated: len=%d", len(*out[0].Live))
	}
	if len(out[0].Message) > 257 {
		t.Errorf("message not truncated: len=%d", len(out[0].Message))
	}
}

func TestTruncator_ShortValues_Unchanged(t *testing.T) {
	tr := NewTruncator(DefaultTruncateConfig())
	out := tr.Apply(truncResults())

	if *out[1].Declared != "short" {
		t.Errorf("expected 'short', got %q", *out[1].Declared)
	}
	if out[1].Message != "ok" {
		t.Errorf("expected 'ok', got %q", out[1].Message)
	}
}

func TestTruncator_NilPointers(t *testing.T) {
	tr := NewTruncator(DefaultTruncateConfig())
	out := tr.Apply([]DriftResult{{Key: "x", Declared: nil, Live: nil}})
	if out[0].Declared != nil || out[0].Live != nil {
		t.Error("nil pointers should remain nil")
	}
}

func TestTruncator_DefaultsOnZeroConfig(t *testing.T) {
	tr := NewTruncator(TruncateConfig{})
	if tr.cfg.MaxValueLen != 128 {
		t.Errorf("expected default MaxValueLen=128, got %d", tr.cfg.MaxValueLen)
	}
	if tr.cfg.MaxMessageLen != 256 {
		t.Errorf("expected default MaxMessageLen=256, got %d", tr.cfg.MaxMessageLen)
	}
}

func TestTruncator_NoMutation(t *testing.T) {
	original := truncResults()
	origLen := len(*original[0].Declared)
	tr := NewTruncator(DefaultTruncateConfig())
	tr.Apply(original)
	if len(*original[0].Declared) != origLen {
		t.Error("Apply should not mutate original results")
	}
}
