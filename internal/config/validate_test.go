package config

import (
	"strings"
	"testing"
)

func TestValidateRuleRaw_Build_MissingPrefix(t *testing.T) {
	r := ValidateRuleRaw{Pattern: `^\d+$`}
	_, err := r.Build()
	if err == nil || !strings.Contains(err.Error(), "prefix") {
		t.Errorf("expected prefix error, got %v", err)
	}
}

func TestValidateRuleRaw_Build_MissingPattern(t *testing.T) {
	r := ValidateRuleRaw{Prefix: "app."}
	_, err := r.Build()
	if err == nil || !strings.Contains(err.Error(), "pattern") {
		t.Errorf("expected pattern error, got %v", err)
	}
}

func TestValidateRuleRaw_Build_InvalidPattern(t *testing.T) {
	r := ValidateRuleRaw{Prefix: "app.", Pattern: `[invalid`}
	_, err := r.Build()
	if err == nil || !strings.Contains(err.Error(), "invalid pattern") {
		t.Errorf("expected invalid pattern error, got %v", err)
	}
}

func TestValidateRuleRaw_Build_Valid(t *testing.T) {
	r := ValidateRuleRaw{Prefix: "db.", Pattern: `^\d+$`, Message: "must be numeric"}
	rule, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Prefix != "db." {
		t.Errorf("unexpected prefix: %s", rule.Prefix)
	}
	if rule.Message != "must be numeric" {
		t.Errorf("unexpected message: %s", rule.Message)
	}
	if rule.Pattern == nil {
		t.Error("expected compiled pattern")
	}
}

func TestBuildValidateRules_Valid(t *testing.T) {
	raws := []ValidateRuleRaw{
		{Prefix: "db.", Pattern: `^\d+$`},
		{Prefix: "app.", Pattern: `^[a-z]+$`, Message: "lowercase only"},
	}
	rules, err := BuildValidateRules(raws)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestBuildValidateRules_PropagatesError(t *testing.T) {
	raws := []ValidateRuleRaw{
		{Prefix: "db.", Pattern: `^\d+$`},
		{Prefix: "app.", Pattern: `[bad`},
	}
	_, err := BuildValidateRules(raws)
	if err == nil {
		t.Error("expected error from invalid pattern")
	}
}
