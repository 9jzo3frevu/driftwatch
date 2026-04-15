package drift

import (
	"testing"
)

func TestDetect_NoDrift(t *testing.T) {
	expected := State{"replicas": 3, "image": "nginx:1.25", "port": 80}
	actual := State{"replicas": 3, "image": "nginx:1.25", "port": 80}

	result := Detect("web", expected, actual)

	if result.Drifted {
		t.Errorf("expected no drift, got %d change(s)", len(result.Changes))
	}
	if len(result.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(result.Changes))
	}
}

func TestDetect_ModifiedValue(t *testing.T) {
	expected := State{"replicas": 3}
	actual := State{"replicas": 5}

	result := Detect("api", expected, actual)

	if !result.Drifted {
		t.Fatal("expected drift to be detected")
	}
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Kind != ChangeModified {
		t.Errorf("expected kind %q, got %q", ChangeModified, result.Changes[0].Kind)
	}
}

func TestDetect_RemovedKey(t *testing.T) {
	expected := State{"image": "redis:7", "port": 6379}
	actual := State{"image": "redis:7"}

	result := Detect("cache", expected, actual)

	if !result.Drifted {
		t.Fatal("expected drift to be detected")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "port" && c.Kind == ChangeRemoved {
			found = true
		}
	}
	if !found {
		t.Error("expected a removed change for key 'port'")
	}
}

func TestDetect_AddedKey(t *testing.T) {
	expected := State{"image": "postgres:15"}
	actual := State{"image": "postgres:15", "debug": true}

	result := Detect("db", expected, actual)

	if !result.Drifted {
		t.Fatal("expected drift to be detected")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "debug" && c.Kind == ChangeAdded {
			found = true
		}
	}
	if !found {
		t.Error("expected an added change for key 'debug'")
	}
}

func TestDriftResult_Summary(t *testing.T) {
	noDrift := DriftResult{ServiceName: "svc", Drifted: false}
	if got := noDrift.Summary(); got != "[OK] svc: no drift detected" {
		t.Errorf("unexpected summary: %s", got)
	}

	withDrift := DriftResult{
		ServiceName: "svc",
		Drifted:     true,
		Changes:     []Change{{Key: "x", Kind: ChangeModified}},
	}
	if got := withDrift.Summary(); got != "[DRIFT] svc: 1 change(s) detected" {
		t.Errorf("unexpected summary: %s", got)
	}
}
