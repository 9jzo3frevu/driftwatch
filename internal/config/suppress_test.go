package config

import (
	"testing"
	"time"
)

func TestSuppressionRaw_Build_MissingKey(t *testing.T) {
	r := SuppressionRaw{Expires: "1h"}
	_, err := r.Build()
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestSuppressionRaw_Build_MissingExpires(t *testing.T) {
	r := SuppressionRaw{Key: "db.host"}
	_, err := r.Build()
	if err == nil {
		t.Error("expected error for missing expires")
	}
}

func TestSuppressionRaw_Build_InvalidDuration(t *testing.T) {
	r := SuppressionRaw{Key: "db.host", Expires: "notaduration"}
	_, err := r.Build()
	if err == nil {
		t.Error("expected error for invalid duration")
	}
}

func TestSuppressionRaw_Build_Valid(t *testing.T) {
	before := time.Now()
	r := SuppressionRaw{Key: "app.port", Expires: "2h"}
	rule, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Key != "app.port" {
		t.Errorf("expected key app.port, got %s", rule.Key)
	}
	expected := before.Add(2 * time.Hour)
	if rule.ExpiresAt.Before(expected.Add(-5*time.Second)) || rule.ExpiresAt.After(expected.Add(5*time.Second)) {
		t.Errorf("ExpiresAt out of expected range")
	}
}

func TestBuildSuppressions_Valid(t *testing.T) {
	raw := []SuppressionRaw{
		{Key: "a", Expires: "1h"},
		{Key: "b", Expires: "30m"},
	}
	rules, err := BuildSuppressions(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestBuildSuppressions_Error(t *testing.T) {
	raw := []SuppressionRaw{
		{Key: "a", Expires: "1h"},
		{Key: "", Expires: "30m"},
	}
	_, err := BuildSuppressions(raw)
	if err == nil {
		t.Error("expected error from invalid rule")
	}
}
