package drift

import (
	"testing"
)

func enrichResults() []DriftResult {
	return []DriftResult{
		{Key: "db.host", Declared: ptrStr("localhost"), Live: ptrStr("prod-db")},
		{Key: "cache.ttl", Declared: ptrStr("60"), Live: ptrStr("120")},
		{Key: "db.port", Declared: ptrStr("5432"), Live: ptrStr("5433")},
	}
}

func TestEnricher_NoRules(t *testing.T) {
	e := NewEnricher(EnrichConfig{})
	out := e.Enrich(enrichResults())
	for _, r := range out {
		if len(r.Metadata) != 0 {
			t.Errorf("expected no metadata for key %s", r.Key)
		}
	}
}

func TestEnricher_PrefixMatch(t *testing.T) {
	e := NewEnricher(EnrichConfig{
		Rules: []EnrichRule{
			{Prefix: "db.", Metadata: map[string]string{"team": "data", "tier": "critical"}},
		},
	})
	out := e.Enrich(enrichResults())
	for _, r := range out {
		if r.Key == "db.host" || r.Key == "db.port" {
			if r.Metadata["team"] != "data" {
				t.Errorf("expected team=data for %s", r.Key)
			}
			if r.Metadata["tier"] != "critical" {
				t.Errorf("expected tier=critical for %s", r.Key)
			}
		}
		if r.Key == "cache.ttl" && len(r.Metadata) != 0 {
			t.Errorf("cache.ttl should have no metadata")
		}
	}
}

func TestEnricher_MultipleRules(t *testing.T) {
	e := NewEnricher(EnrichConfig{
		Rules: []EnrichRule{
			{Prefix: "db.", Metadata: map[string]string{"team": "data"}},
			{Prefix: "cache.", Metadata: map[string]string{"team": "platform"}},
		},
	})
	out := e.Enrich(enrichResults())
	teams := map[string]string{}
	for _, r := range out {
		if r.Metadata != nil {
			teams[r.Key] = r.Metadata["team"]
		}
	}
	if teams["db.host"] != "data" {
		t.Errorf("expected data team for db.host")
	}
	if teams["cache.ttl"] != "platform" {
		t.Errorf("expected platform team for cache.ttl")
	}
}

func TestEnricher_NoMutation(t *testing.T) {
	original := enrichResults()
	e := NewEnricher(EnrichConfig{
		Rules: []EnrichRule{
			{Prefix: "db.", Metadata: map[string]string{"owner": "dba"}},
		},
	})
	e.Enrich(original)
	for _, r := range original {
		if len(r.Metadata) != 0 {
			t.Errorf("original results should not be mutated")
		}
	}
}
