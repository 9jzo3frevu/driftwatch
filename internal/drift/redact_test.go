package drift

import (
	"testing"
)

func redactResults() []DriftResult {
	return []DriftResult{
		{Key: "app.name", Declared: ptrStr("myapp"), Live: ptrStr("myapp")},
		{Key: "db.password", Declared: ptrStr("secret123"), Live: ptrStr("different")},
		{Key: "api_key", Declared: ptrStr("abc"), Live: nil},
		{Key: "auth.token", Declared: nil, Live: ptrStr("tok")},
	}
}

func TestRedactor_NonSensitiveUnchanged(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	results := redactResults()
	out := r.Redact(results)
	if *out[0].Declared != "myapp" || *out[0].Live != "myapp" {
		t.Errorf("expected app.name unchanged, got %v / %v", out[0].Declared, out[0].Live)
	}
}

func TestRedactor_MasksPassword(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	out := r.Redact(redactResults())
	if *out[1].Declared != "***REDACTED***" {
		t.Errorf("expected declared masked, got %v", *out[1].Declared)
	}
	if *out[1].Live != "***REDACTED***" {
		t.Errorf("expected live masked, got %v", *out[1].Live)
	}
}

func TestRedactor_MasksAPIKey(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	out := r.Redact(redactResults())
	if out[2].Declared == nil || *out[2].Declared != "***REDACTED***" {
		t.Errorf("expected api_key declared masked")
	}
	if out[2].Live != nil {
		t.Errorf("expected nil live unchanged")
	}
}

func TestRedactor_MasksToken(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	out := r.Redact(redactResults())
	if out[3].Live == nil || *out[3].Live != "***REDACTED***" {
		t.Errorf("expected token live masked")
	}
}

func TestRedactor_CustomMask(t *testing.T) {
	cfg := RedactConfig{Patterns: []string{"secret"}, Mask: "[hidden]"}
	r := NewRedactor(cfg)
	results := []DriftResult{{Key: "my.secret", Declared: ptrStr("val"), Live: ptrStr("val2")}}
	out := r.Redact(results)
	if *out[0].Declared != "[hidden]" || *out[0].Live != "[hidden]" {
		t.Errorf("expected custom mask, got %v / %v", *out[0].Declared, *out[0].Live)
	}
}

func TestRedactor_NoMutation(t *testing.T) {
	r := NewRedactor(DefaultRedactConfig())
	original := redactResults()
	r.Redact(original)
	if *original[1].Declared != "secret123" {
		t.Errorf("original results should not be mutated")
	}
}
