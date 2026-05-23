package drift

import (
	"regexp"
	"testing"
)

func validateResults() []DriftResult {
	return []DriftResult{
		{Key: "db.port", Declared: "5432", Live: ptrStr("5432"), Drifted: false},
		{Key: "db.host", Declared: "localhost", Live: ptrStr("192.168.1.1"), Drifted: true},
		{Key: "app.version", Declared: "1.0.0", Live: ptrStr("not-a-version"), Drifted: true},
		{Key: "app.name", Declared: "svc", Live: nil, Drifted: true},
	}
}

func TestValidator_NoRules(t *testing.T) {
	v := NewValidator(ValidateConfig{})
	res := v.Validate(validateResults())
	for _, r := range res {
		if len(r.Annotations) != 0 {
			t.Errorf("expected no annotations, got %v", r.Annotations)
		}
	}
}

func TestValidator_MatchingPattern(t *testing.T) {
	v := NewValidator(ValidateConfig{
		Rules: []ValidateRule{
			{Prefix: "app.version", Pattern: regexp.MustCompile(`^\d+\.\d+\.\d+$`), Message: "must be semver"},
		},
	})
	res := v.Validate(validateResults())
	var found bool
	for _, r := range res {
		if r.Key == "app.version" {
			if len(r.Annotations) == 0 {
				t.Fatal("expected validation annotation on app.version")
			}
			if r.Annotations[0] != "validation_error:must be semver" {
				t.Errorf("unexpected annotation: %s", r.Annotations[0])
			}
			found = true
		}
	}
	if !found {
		t.Fatal("app.version result not found")
	}
}

func TestValidator_NilLive_Skipped(t *testing.T) {
	v := NewValidator(ValidateConfig{
		Rules: []ValidateRule{
			{Prefix: "app.name", Pattern: regexp.MustCompile(`^[a-z]+$`)},
		},
	})
	res := v.Validate(validateResults())
	for _, r := range res {
		if r.Key == "app.name" && len(r.Annotations) != 0 {
			t.Errorf("nil live should not produce annotation")
		}
	}
}

func TestValidator_NoMutation(t *testing.T) {
	orig := validateResults()
	v := NewValidator(ValidateConfig{
		Rules: []ValidateRule{
			{Prefix: "db.", Pattern: regexp.MustCompile(`^\d+$`)},
		},
	})
	v.Validate(orig)
	if len(orig[0].Annotations) != 0 {
		t.Error("original results should not be mutated")
	}
}

func TestValidator_DefaultMessage(t *testing.T) {
	v := NewValidator(ValidateConfig{
		Rules: []ValidateRule{
			{Prefix: "db.host", Pattern: regexp.MustCompile(`^localhost$`)},
		},
	})
	res := v.Validate(validateResults())
	for _, r := range res {
		if r.Key == "db.host" {
			if len(r.Annotations) == 0 {
				t.Fatal("expected annotation")
			}
			if r.Annotations[0] == "" {
				t.Error("expected non-empty default message")
			}
		}
	}
}
