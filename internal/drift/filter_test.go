package drift

import (
	"testing"
)

func filterResults() []DriftResult {
	return []DriftResult{
		{Key: "replicas", Declared: ptrStr("3"), Actual: ptrStr("2")},
		{Key: "image", Declared: ptrStr("app:1.0"), Actual: ptrStr("app:2.0")},
		{Key: "timeout", Declared: ptrStr("30s"), Actual: nil},
		{Key: "debug", Declared: nil, Actual: ptrStr("true")},
	}
}

func TestFilter_NoConstraints(t *testing.T) {
	f := &Filter{}
	out := f.Apply(filterResults())
	if len(out) != 4 {
		t.Fatalf("expected 4 results, got %d", len(out))
	}
}

func TestFilter_ExcludeKey(t *testing.T) {
	f := &Filter{ExcludeKeys: []string{"debug"}}
	out := f.Apply(filterResults())
	for _, r := range out {
		if r.Key == "debug" {
			t.Fatal("excluded key 'debug' should not appear")
		}
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestFilter_IncludeKeys(t *testing.T) {
	f := &Filter{IncludeKeys: []string{"replicas", "image"}}
	out := f.Apply(filterResults())
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestFilter_MinSeverity_FiltersLow(t *testing.T) {
	// replicas diff => low, removed key => higher
	// Use high min to filter most out
	f := &Filter{MinSeverity: "high"}
	out := f.Apply(filterResults())
	for _, r := range out {
		sev := severityFor(r)
		if severityRank[sev] < severityRank["high"] {
			t.Fatalf("result %q has severity %q below minimum", r.Key, sev)
		}
	}
}

func TestFilter_CaseInsensitiveKey(t *testing.T) {
	f := &Filter{ExcludeKeys: []string{"DEBUG"}}
	out := f.Apply(filterResults())
	for _, r := range out {
		if r.Key == "debug" {
			t.Fatal("case-insensitive exclude failed")
		}
	}
}
