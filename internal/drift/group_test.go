package drift

import (
	"testing"
)

func groupResults() []Result {
	return []Result{
		{Key: "host", Service: "api", Declared: ptrStr("a"), Live: ptrStr("b")},
		{Key: "port", Service: "api", Declared: ptrStr("80"), Live: ptrStr("90")},
		{Key: "timeout", Service: "worker", Declared: ptrStr("30"), Live: nil},
	}
}

func TestGroupResults_ByService(t *testing.T) {
	groups := GroupResults(groupResults(), GroupByService)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "api" {
		t.Errorf("expected first group 'api', got %q", groups[0].Name)
	}
	if len(groups[0].Results) != 2 {
		t.Errorf("expected 2 results in api group, got %d", len(groups[0].Results))
	}
}

func TestGroupResults_BySeverity(t *testing.T) {
	groups := GroupResults(groupResults(), GroupBySeverity)
	if len(groups) == 0 {
		t.Fatal("expected at least one group")
	}
	for _, g := range groups {
		if g.Name == "" {
			t.Error("group name should not be empty")
		}
	}
}

func TestGroupResults_ByKey(t *testing.T) {
	groups := GroupResults(groupResults(), GroupByKey)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
}

func TestGroupResults_Empty(t *testing.T) {
	groups := GroupResults([]Result{}, GroupByService)
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestGroupResults_NoService(t *testing.T) {
	results := []Result{
		{Key: "x", Declared: ptrStr("1"), Live: ptrStr("2")},
	}
	groups := GroupResults(results, GroupByService)
	if groups[0].Name != "default" {
		t.Errorf("expected 'default', got %q", groups[0].Name)
	}
}
