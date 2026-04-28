package config

import (
	"testing"
)

func TestLabelRuleRaw_Build_MissingPrefix(t *testing.T) {
	r := LabelRuleRaw{Labels: map[string]string{"team": "data"}}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for missing prefix")
	}
}

func TestLabelRuleRaw_Build_NoLabels(t *testing.T) {
	r := LabelRuleRaw{Prefix: "db."}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for empty labels map")
	}
}

func TestLabelRuleRaw_Build_Valid(t *testing.T) {
	r := LabelRuleRaw{
		Prefix: "db.",
		Labels: map[string]string{"team": "data", "tier": "backend"},
	}
	rule, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Prefix != "db." {
		t.Errorf("expected prefix db., got %q", rule.Prefix)
	}
	if rule.Labels["team"] != "data" {
		t.Errorf("expected label team=data, got %v", rule.Labels)
	}
}

func TestBuildLabelRules_Valid(t *testing.T) {
	raw := []LabelRuleRaw{
		{Prefix: "db.", Labels: map[string]string{"team": "data"}},
		{Prefix: "cache.", Labels: map[string]string{"team": "platform"}},
	}
	rules, err := BuildLabelRules(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestBuildLabelRules_PropagatesError(t *testing.T) {
	raw := []LabelRuleRaw{
		{Prefix: "db.", Labels: map[string]string{"team": "data"}},
		{Prefix: "", Labels: map[string]string{"team": "platform"}},
	}
	_, err := BuildLabelRules(raw)
	if err == nil {
		t.Fatal("expected error to propagate from invalid rule")
	}
}
