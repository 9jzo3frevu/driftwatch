package drift

import (
	"testing"
)

func TestDiff_NoDrift(t *testing.T) {
	declared := map[string]string{"a": "1", "b": "2"}
	live := map[string]string{"a": "1", "b": "2"}
	results := Diff(declared, live, DiffModeExact)
	if len(results) != 0 {
		t.Fatalf("expected no drift, got %d results", len(results))
	}
}

func TestDiff_ModifiedValue(t *testing.T) {
	declared := map[string]string{"a": "1"}
	live := map[string]string{"a": "2"}
	results := Diff(declared, live, DiffModeExact)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "a" || !results[0].Drifted {
		t.Errorf("unexpected result: %+v", results[0])
	}
}

func TestDiff_MissingLiveKey(t *testing.T) {
	declared := map[string]string{"a": "1", "b": "2"}
	live := map[string]string{"a": "1"}
	results := Diff(declared, live, DiffModeExact)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "b" {
		t.Errorf("expected missing key 'b', got %q", results[0].Key)
	}
	if results[0].Actual != nil {
		t.Error("expected Actual to be nil for missing key")
	}
}

func TestDiff_ExtraLiveKey_ExactMode(t *testing.T) {
	declared := map[string]string{"a": "1"}
	live := map[string]string{"a": "1", "b": "2"}
	results := Diff(declared, live, DiffModeExact)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "b" {
		t.Errorf("expected extra key 'b', got %q", results[0].Key)
	}
}

func TestDiff_ExtraLiveKey_SubsetMode(t *testing.T) {
	declared := map[string]string{"a": "1"}
	live := map[string]string{"a": "1", "b": "2"}
	results := Diff(declared, live, DiffModeSubset)
	if len(results) != 0 {
		t.Fatalf("expected no drift in subset mode, got %d results", len(results))
	}
}

func TestDiff_MultipleIssues(t *testing.T) {
	declared := map[string]string{"a": "1", "b": "2", "c": "3"}
	live := map[string]string{"a": "changed", "d": "extra"}
	results := Diff(declared, live, DiffModeExact)
	// a modified, b missing, c missing, d extra = 4
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
}
