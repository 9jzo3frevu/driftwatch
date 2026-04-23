package drift

import (
	"testing"
)

func normalizeResults() []DriftResult {
	return []DriftResult{
		{Key: "host", Declared: ptrStr("  localhost  "), Live: ptrStr("  localhost  "), Drift: false},
		{Key: "env", Declared: ptrStr(`"production"`), Live: ptrStr(`"staging"`), Drift: true},
		{Key: "mode", Declared: ptrStr("DEBUG"), Live: ptrStr("debug"), Drift: true},
		{Key: "missing", Declared: ptrStr("value"), Live: nil, Drift: true},
	}
}

func TestNormalizer_TrimSpace(t *testing.T) {
	cfg := NormalizeConfig{TrimSpace: true}
	n := NewNormalizer(cfg)
	results := n.Normalize(normalizeResults())

	if *results[0].Declared != "localhost" {
		t.Errorf("expected trimmed declared, got %q", *results[0].Declared)
	}
	if *results[0].Live != "localhost" {
		t.Errorf("expected trimmed live, got %q", *results[0].Live)
	}
}

func TestNormalizer_StripQuotes(t *testing.T) {
	cfg := NormalizeConfig{StripQuotes: true}
	n := NewNormalizer(cfg)
	results := n.Normalize(normalizeResults())

	if *results[1].Declared != "production" {
		t.Errorf("expected unquoted declared, got %q", *results[1].Declared)
	}
	if *results[1].Live != "staging" {
		t.Errorf("expected unquoted live, got %q", *results[1].Live)
	}
}

func TestNormalizer_Lowercase(t *testing.T) {
	cfg := NormalizeConfig{Lowercase: true}
	n := NewNormalizer(cfg)
	results := n.Normalize(normalizeResults())

	if *results[2].Declared != "debug" {
		t.Errorf("expected lowercased declared, got %q", *results[2].Declared)
	}
	if *results[2].Live != "debug" {
		t.Errorf("expected lowercased live, got %q", *results[2].Live)
	}
}

func TestNormalizer_NilLive_Preserved(t *testing.T) {
	cfg := DefaultNormalizeConfig()
	n := NewNormalizer(cfg)
	results := n.Normalize(normalizeResults())

	if results[3].Live != nil {
		t.Errorf("expected nil live to remain nil")
	}
	if *results[3].Declared != "value" {
		t.Errorf("expected declared unchanged, got %q", *results[3].Declared)
	}
}

func TestNormalizer_NoMutation(t *testing.T) {
	original := normalizeResults()
	cfg := NormalizeConfig{TrimSpace: true, Lowercase: true, StripQuotes: true}
	n := NewNormalizer(cfg)
	_ = n.Normalize(original)

	if *original[0].Declared != "  localhost  " {
		t.Errorf("original results mutated unexpectedly")
	}
}
