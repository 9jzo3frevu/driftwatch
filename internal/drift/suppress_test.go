package drift

import (
	"testing"
	"time"
)

func suppressResults() []DriftResult {
	return []DriftResult{
		{Key: "db.host", Declared: ptrStr("localhost"), Actual: ptrStr("prod-db")},
		{Key: "app.port", Declared: ptrStr("8080"), Actual: ptrStr("9090")},
		{Key: "cache.ttl", Declared: ptrStr("300"), Actual: ptrStr("600")},
	}
}

func TestSuppressor_NoneExpired(t *testing.T) {
	now := time.Now()
	rules := []SuppressionRule{
		{Key: "db.host", ExpiresAt: now.Add(1 * time.Hour)},
	}
	s := NewSuppressor(rules)
	out := s.Apply(suppressResults(), now)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.Key == "db.host" {
			t.Error("suppressed key should not appear")
		}
	}
}

func TestSuppressor_Expired(t *testing.T) {
	now := time.Now()
	rules := []SuppressionRule{
		{Key: "db.host", ExpiresAt: now.Add(-1 * time.Hour)},
	}
	s := NewSuppressor(rules)
	out := s.Apply(suppressResults(), now)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestSuppressor_NoRules(t *testing.T) {
	s := NewSuppressor(nil)
	out := s.Apply(suppressResults(), time.Now())
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestSuppressor_IsSuppressed(t *testing.T) {
	now := time.Now()
	rules := []SuppressionRule{
		{Key: "app.port", ExpiresAt: now.Add(10 * time.Minute)},
	}
	s := NewSuppressor(rules)
	if !s.IsSuppressed("app.port", now) {
		t.Error("expected app.port to be suppressed")
	}
	if s.IsSuppressed("db.host", now) {
		t.Error("expected db.host to not be suppressed")
	}
}
