package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func pinResults() []DriftResult {
	return []DriftResult{
		{Key: "db.host", Service: "api", Drifted: true},
		{Key: "db.port", Service: "api", Drifted: true},
		{Key: "cache.ttl", Service: "worker", Drifted: true},
	}
}

func TestPinnedKey_IsExpired_NoExpiry(t *testing.T) {
	p := PinnedKey{Key: "k", PinnedAt: time.Now()}
	if p.IsExpired(time.Now().Add(24 * time.Hour)) {
		t.Fatal("expected pin without expiry to never expire")
	}
}

func TestPinnedKey_IsExpired_Past(t *testing.T) {
	p := PinnedKey{
		Key:       "k",
		PinnedAt:  time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if !p.IsExpired(time.Now()) {
		t.Fatal("expected expired pin to report true")
	}
}

func TestPinStore_SaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	store := NewPinStore(path)

	now := time.Now().UTC().Truncate(time.Second)
	pins := []PinnedKey{
		{Key: "db.host", Service: "api", PinnedAt: now, Reason: "manual"},
		{Key: "cache.ttl", Service: "worker", PinnedAt: now, ExpiresAt: now.Add(time.Hour)},
	}
	if err := store.Save(pins); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := store.Load(now)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 pins, got %d", len(loaded))
	}
}

func TestPinStore_Load_FiltersExpired(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	now := time.Now().UTC()
	pins := []PinnedKey{
		{Key: "a", Service: "svc", PinnedAt: now, ExpiresAt: now.Add(-time.Minute)},
		{Key: "b", Service: "svc", PinnedAt: now},
	}
	data, _ := json.Marshal(pins)
	_ = os.WriteFile(path, data, 0o644)

	store := NewPinStore(path)
	loaded, err := store.Load(now)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Key != "b" {
		t.Fatalf("expected only non-expired pin, got %+v", loaded)
	}
}

func TestPinStore_Load_MissingFile(t *testing.T) {
	store := NewPinStore("/nonexistent/pins.json")
	pins, err := store.Load(time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pins) != 0 {
		t.Fatalf("expected empty slice, got %d", len(pins))
	}
}

func TestApplyPins_RemovesPinned(t *testing.T) {
	results := pinResults()
	pins := []PinnedKey{
		{Key: "db.host", Service: "api"},
	}
	out := ApplyPins(results, pins)
	if len(out) != 2 {
		t.Fatalf("expected 2 results after pin, got %d", len(out))
	}
	for _, r := range out {
		if r.Key == "db.host" && r.Service == "api" {
			t.Fatal("pinned key should have been removed")
		}
	}
}

func TestApplyPins_NoPins_ReturnsAll(t *testing.T) {
	results := pinResults()
	out := ApplyPins(results, nil)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}
