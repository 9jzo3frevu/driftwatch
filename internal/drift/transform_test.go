package drift

import (
	"testing"
)

func transformResults() []DriftResult {
	return []DriftResult{
		{Key: "app.name", Declared: ptrStr("myapp"), Live: ptrStr("myapp")},
		{Key: "app.env", Declared: ptrStr("production"), Live: ptrStr("staging")},
		{Key: "db.host", Declared: ptrStr("localhost"), Live: nil},
	}
}

func TestTransformer_NoRules(t *testing.T) {
	tr := NewTransformer(TransformConfig{})
	results := transformResults()
	out := tr.Apply(results)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
	if *out[0].Declared != "myapp" {
		t.Errorf("expected unchanged declared, got %q", *out[0].Declared)
	}
}

func TestTransformer_PrependOp(t *testing.T) {
	cfg := TransformConfig{
		Rules: []TransformRule{
			{Prefix: "app.", Op: TransformOpPrefix, Arg: "svc-"},
		},
	}
	tr := NewTransformer(cfg)
	out := tr.Apply(transformResults())
	if *out[0].Declared != "svc-myapp" {
		t.Errorf("expected svc-myapp, got %q", *out[0].Declared)
	}
	if *out[0].Live != "svc-myapp" {
		t.Errorf("expected svc-myapp live, got %q", *out[0].Live)
	}
	// db.host should be unaffected
	if *out[2].Declared != "localhost" {
		t.Errorf("expected localhost unchanged, got %q", *out[2].Declared)
	}
}

func TestTransformer_UpperOp(t *testing.T) {
	cfg := TransformConfig{
		Rules: []TransformRule{
			{Prefix: "app.env", Op: TransformOpUpper},
		},
	}
	tr := NewTransformer(cfg)
	out := tr.Apply(transformResults())
	if *out[1].Declared != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", *out[1].Declared)
	}
	if *out[1].Live != "STAGING" {
		t.Errorf("expected STAGING, got %q", *out[1].Live)
	}
}

func TestTransformer_NilLive_Preserved(t *testing.T) {
	cfg := TransformConfig{
		Rules: []TransformRule{
			{Prefix: "db.", Op: TransformOpSuffix, Arg: ":5432"},
		},
	}
	tr := NewTransformer(cfg)
	out := tr.Apply(transformResults())
	if out[2].Live != nil {
		t.Errorf("expected nil live to remain nil")
	}
	if *out[2].Declared != "localhost:5432" {
		t.Errorf("expected localhost:5432, got %q", *out[2].Declared)
	}
}

func TestTransformer_NoMutation(t *testing.T) {
	cfg := TransformConfig{
		Rules: []TransformRule{
			{Prefix: "app.", Op: TransformOpReplace, Arg: "replaced"},
		},
	}
	tr := NewTransformer(cfg)
	originals := transformResults()
	tr.Apply(originals)
	if *originals[0].Declared != "myapp" {
		t.Errorf("original result was mutated")
	}
}
