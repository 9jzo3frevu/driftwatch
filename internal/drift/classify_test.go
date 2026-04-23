package drift

import (
	"testing"
)

func classifyResults() []DriftResult {
	return []DriftResult{
		{Key: "auth.token", Declared: ptrStr("abc"), Live: ptrStr("xyz")},
		{Key: "db.host", Declared: ptrStr("localhost"), Live: ptrStr("10.0.0.1")},
		{Key: "server.port", Declared: ptrStr("8080"), Live: ptrStr("9090")},
		{Key: "app.name", Declared: ptrStr("driftwatch"), Live: ptrStr("other")},
		{Key: "tls.cert.path", Declared: ptrStr("/etc/certs/a"), Live: ptrStr("/etc/certs/b")},
	}
}

func TestClassifier_DefaultConfig_Security(t *testing.T) {
	c := NewClassifier(DefaultClassifyConfig())
	results := c.Classify(classifyResults())
	if results[0].Category != "security" {
		t.Errorf("expected security for auth.token, got %q", results[0].Category)
	}
}

func TestClassifier_DefaultConfig_Storage(t *testing.T) {
	c := NewClassifier(DefaultClassifyConfig())
	results := c.Classify(classifyResults())
	if results[1].Category != "storage" {
		t.Errorf("expected storage for db.host, got %q", results[1].Category)
	}
}

func TestClassifier_DefaultConfig_Network(t *testing.T) {
	c := NewClassifier(DefaultClassifyConfig())
	results := c.Classify(classifyResults())
	if results[2].Category != "network" {
		t.Errorf("expected network for server.port, got %q", results[2].Category)
	}
}

func TestClassifier_DefaultConfig_General(t *testing.T) {
	c := NewClassifier(DefaultClassifyConfig())
	results := c.Classify(classifyResults())
	if results[3].Category != "general" {
		t.Errorf("expected general for app.name, got %q", results[3].Category)
	}
}

func TestClassifier_DefaultConfig_CertIsSecurity(t *testing.T) {
	c := NewClassifier(DefaultClassifyConfig())
	results := c.Classify(classifyResults())
	if results[4].Category != "security" {
		t.Errorf("expected security for tls.cert.path, got %q", results[4].Category)
	}
}

func TestClassifier_NoMutation(t *testing.T) {
	c := NewClassifier(DefaultClassifyConfig())
	original := classifyResults()
	c.Classify(original)
	for _, r := range original {
		if r.Category != "" {
			t.Errorf("original results should not be mutated, got category %q", r.Category)
		}
	}
}

func TestClassifier_EmptyResults(t *testing.T) {
	c := NewClassifier(DefaultClassifyConfig())
	results := c.Classify([]DriftResult{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestClassifier_CustomConfig(t *testing.T) {
	cfg := ClassifyConfig{
		CategoryRules:   map[string][]string{"infra": {"app"}},
		DefaultCategory: "unknown",
	}
	c := NewClassifier(cfg)
	results := c.Classify([]DriftResult{
		{Key: "app.name", Declared: ptrStr("a"), Live: ptrStr("b")},
		{Key: "other.key", Declared: ptrStr("x"), Live: ptrStr("y")},
	})
	if results[0].Category != "infra" {
		t.Errorf("expected infra, got %q", results[0].Category)
	}
	if results[1].Category != "unknown" {
		t.Errorf("expected unknown, got %q", results[1].Category)
	}
}
