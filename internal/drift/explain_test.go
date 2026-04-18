package drift

import (
	"testing"
)

func explainResults() []DriftResult {
	decl := "v1"
	live := "v2"
	return []DriftResult{
		{Key: "app.version", Kind: KindModified, Declared: &decl, Live: &live},
		{Key: "app.timeout", Kind: KindMissing, Declared: &decl, Live: nil},
		{Key: "app.debug", Kind: KindExtra, Declared: nil, Live: &live},
	}
}

func TestExplainer_NoHints(t *testing.T) {
	ex := NewExplainer(nil)
	results := explainResults()
	exps := ex.Explain(results)
	if len(exps) != 3 {
		t.Fatalf("expected 3 explanations, got %d", len(exps))
	}
	if exps[0].Hint != "" {
		t.Errorf("expected empty hint, got %q", exps[0].Hint)
	}
}

func TestExplainer_WithHint(t *testing.T) {
	hints := map[string]string{"app.version": "check deploy pipeline"}
	ex := NewExplainer(hints)
	exps := ex.Explain(explainResults())
	if exps[0].Hint != "check deploy pipeline" {
		t.Errorf("unexpected hint: %q", exps[0].Hint)
	}
	if exps[1].Hint != "" {
		t.Errorf("expected empty hint for app.timeout")
	}
}

func TestExplainer_MessageContent(t *testing.T) {
	ex := NewExplainer(nil)
	exps := ex.Explain(explainResults())

	if exps[0].Kind != string(KindModified) {
		t.Errorf("expected modified kind")
	}
	if exps[1].Kind != string(KindMissing) {
		t.Errorf("expected missing kind")
	}
	if exps[2].Kind != string(KindExtra) {
		t.Errorf("expected extra kind")
	}
}

func TestExplainer_EmptyResults(t *testing.T) {
	ex := NewExplainer(nil)
	exps := ex.Explain([]DriftResult{})
	if len(exps) != 0 {
		t.Errorf("expected empty explanations")
	}
}
