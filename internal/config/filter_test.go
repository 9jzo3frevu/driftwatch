package config

import (
	"testing"
)

func TestFilterRaw_Build_Empty(t *testing.T) {
	r := FilterRaw{}
	f, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if len(f.IncludeKeys) != 0 || len(f.ExcludeKeys) != 0 {
		t.Fatal("expected empty key lists")
	}
}

func TestFilterRaw_Build_ValidSeverity(t *testing.T) {
	for _, sev := range []string{"low", "medium", "high", "critical"} {
		r := FilterRaw{MinSeverity: sev}
		f, err := r.Build()
		if err != nil {
			t.Fatalf("severity %q: unexpected error: %v", sev, err)
		}
		if f.MinSeverity != sev {
			t.Fatalf("expected %q, got %q", sev, f.MinSeverity)
		}
	}
}

func TestFilterRaw_Build_InvalidSeverity(t *testing.T) {
	r := FilterRaw{MinSeverity: "extreme"}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for invalid severity")
	}
}

func TestFilterRaw_Build_KeyLists(t *testing.T) {
	r := FilterRaw{
		IncludeKeys: []string{"replicas", "image"},
		ExcludeKeys: []string{"debug"},
	}
	f, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.IncludeKeys) != 2 {
		t.Fatalf("expected 2 include keys, got %d", len(f.IncludeKeys))
	}
	if len(f.ExcludeKeys) != 1 {
		t.Fatalf("expected 1 exclude key, got %d", len(f.ExcludeKeys))
	}
}
