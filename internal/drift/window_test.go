package drift

import (
	"testing"
	"time"
)

func windowResults() []DriftResult {
	return []DriftResult{
		{Key: "app.replicas", Declared: "3", Live: ptrStr("2"), Drifted: true},
		{Key: "app.image", Declared: "v1.0", Live: ptrStr("v1.1"), Drifted: true},
		{Key: "app.port", Declared: "8080", Live: ptrStr("8080"), Drifted: false},
	}
}

func TestWindow_AddAndEntries(t *testing.T) {
	w := NewWindow(WindowConfig{Size: time.Minute, MaxResults: 100})
	w.Add(windowResults())

	entries := w.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestWindow_Results_ReturnsDriftResults(t *testing.T) {
	w := NewWindow(DefaultWindowConfig())
	w.Add(windowResults())

	results := w.Results()
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Key != "app.replicas" {
		t.Errorf("unexpected first key: %s", results[0].Key)
	}
}

func TestWindow_Evicts_ExpiredEntries(t *testing.T) {
	w := NewWindow(WindowConfig{Size: 10 * time.Second, MaxResults: 100})

	base := time.Now()
	w.now = func() time.Time { return base }
	w.Add(windowResults())

	// advance time past the window
	w.now = func() time.Time { return base.Add(15 * time.Second) }
	w.Add([]DriftResult{{Key: "new.key", Declared: "x", Live: ptrStr("y"), Drifted: true}})

	entries := w.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", len(entries))
	}
	if entries[0].Result.Key != "new.key" {
		t.Errorf("expected new.key, got %s", entries[0].Result.Key)
	}
}

func TestWindow_MaxResults_Cap(t *testing.T) {
	w := NewWindow(WindowConfig{Size: time.Minute, MaxResults: 2})
	w.Add(windowResults()) // 3 results, only last 2 kept

	if len(w.Entries()) != 2 {
		t.Fatalf("expected 2 entries due to cap, got %d", len(w.Entries()))
	}
}

func TestWindow_Flush_ClearsEntries(t *testing.T) {
	w := NewWindow(DefaultWindowConfig())
	w.Add(windowResults())
	w.Flush()

	if len(w.Entries()) != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", len(w.Entries()))
	}
}

func TestWindow_DefaultConfig_UsedOnZero(t *testing.T) {
	w := NewWindow(WindowConfig{})
	if w.cfg.Size != DefaultWindowConfig().Size {
		t.Errorf("expected default size %v, got %v", DefaultWindowConfig().Size, w.cfg.Size)
	}
	if w.cfg.MaxResults != DefaultWindowConfig().MaxResults {
		t.Errorf("expected default max results %d, got %d", DefaultWindowConfig().MaxResults, w.cfg.MaxResults)
	}
}
