package drift

import (
	"testing"
)

func labelResults() []DriftResult {
	return []DriftResult{
		{Key: "db.host", Declared: "localhost", Live: ptrStr("remotehost"), Severity: SeverityHigh},
		{Key: "cache.ttl", Declared: "60", Live: ptrStr("120"), Severity: SeverityLow},
		{Key: "db.port", Declared: "5432", Live: ptrStr("5433"), Severity: SeverityMedium},
	}
}

func TestLabeler_NoRules(t *testing.T) {
	l := NewLabeler(LabelConfig{})
	results := labelResults()
	out := l.Apply(results)
	for _, r := range out {
		if len(r.Labels) != 0 {
			t.Errorf("expected no labels, got %v", r.Labels)
		}
	}
}

func TestLabeler_PrefixMatch(t *testing.T) {
	l := NewLabeler(LabelConfig{
		Rules: []LabelRule{
			{Prefix: "db.", Labels: map[string]string{"team": "data", "tier": "backend"}},
		},
	})
	out := l.Apply(labelResults())
	for _, r := range out {
		if r.Key == "db.host" || r.Key == "db.port" {
			if r.Labels["team"] != "data" {
				t.Errorf("key %s: expected label team=data, got %v", r.Key, r.Labels)
			}
			if r.Labels["tier"] != "backend" {
				t.Errorf("key %s: expected label tier=backend, got %v", r.Key, r.Labels)
			}
		}
		if r.Key == "cache.ttl" && len(r.Labels) != 0 {
			t.Errorf("cache.ttl should have no labels, got %v", r.Labels)
		}
	}
}

func TestLabeler_MultipleRulesMatch(t *testing.T) {
	l := NewLabeler(LabelConfig{
		Rules: []LabelRule{
			{Prefix: "db.", Labels: map[string]string{"team": "data"}},
			{Prefix: "db.host", Labels: map[string]string{"critical": "true"}},
		},
	})
	out := l.Apply(labelResults())
	for _, r := range out {
		if r.Key == "db.host" {
			if r.Labels["team"] != "data" || r.Labels["critical"] != "true" {
				t.Errorf("db.host: expected both labels, got %v", r.Labels)
			}
		}
	}
}

func TestLabeler_NoMutation(t *testing.T) {
	l := NewLabeler(LabelConfig{
		Rules: []LabelRule{
			{Prefix: "db.", Labels: map[string]string{"team": "data"}},
		},
	})
	original := labelResults()
	l.Apply(original)
	for _, r := range original {
		if len(r.Labels) != 0 {
			t.Errorf("original results mutated: key %s has labels %v", r.Key, r.Labels)
		}
	}
}
