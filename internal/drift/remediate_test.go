package drift

import (
	"testing"
)

func remediateResults() []DriftResult {
	declared := "v2"
	live := "v1"
	return []DriftResult{
		{Key: "image.tag", Status: StatusModified, Declared: &declared, Live: &live},
		{Key: "replicas", Status: StatusMissing, Declared: strPtr("3"), Live: nil},
		{Key: "env.debug", Status: StatusMatch, Declared: strPtr("false"), Live: strPtr("false")},
	}
}

func TestRemediator_DefaultTemplate(t *testing.T) {
	r := NewRemediator("")
	hints := r.Hints(remediateResults())
	if len(hints) != 2 {
		t.Fatalf("expected 2 hints, got %d", len(hints))
	}
	if hints[0].Key != "image.tag" {
		t.Errorf("expected key image.tag, got %s", hints[0].Key)
	}
	expected := "update image.tag from 'v1' to 'v2'"
	if hints[0].Suggestion != expected {
		t.Errorf("expected %q, got %q", expected, hints[0].Suggestion)
	}
}

func TestRemediator_CustomTemplate(t *testing.T) {
	r := NewRemediator("set {key}={declared}")
	hints := r.Hints(remediateResults())
	if len(hints) != 2 {
		t.Fatalf("expected 2 hints, got %d", len(hints))
	}
	if hints[1].Suggestion != "set replicas=3" {
		t.Errorf("unexpected suggestion: %s", hints[1].Suggestion)
	}
}

func TestRemediator_NoDrift(t *testing.T) {
	r := NewRemediator("")
	results := []DriftResult{
		{Key: "a", Status: StatusMatch, Declared: strPtr("x"), Live: strPtr("x")},
	}
	hints := r.Hints(results)
	if len(hints) != 0 {
		t.Errorf("expected no hints, got %d", len(hints))
	}
}

func TestRemediator_NilLive(t *testing.T) {
	r := NewRemediator("")
	results := []DriftResult{
		{Key: "port", Status: StatusMissing, Declared: strPtr("8080"), Live: nil},
	}
	hints := r.Hints(results)
	if len(hints) != 1 {
		t.Fatalf("expected 1 hint")
	}
	if hints[0].Suggestion == "" {
		t.Error("suggestion should not be empty")
	}
}
