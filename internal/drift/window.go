package drift

import (
	"time"
)

// WindowConfig controls the sliding time window used to aggregate drift results.
type WindowConfig struct {
	// Size is the duration of the window.
	Size time.Duration
	// MaxResults caps the number of results retained within the window.
	MaxResults int
}

// DefaultWindowConfig returns sensible defaults for a sliding window.
func DefaultWindowConfig() WindowConfig {
	return WindowConfig{
		Size:       5 * time.Minute,
		MaxResults: 200,
	}
}

// WindowEntry holds a single drift result with the time it was recorded.
type WindowEntry struct {
	RecordedAt time.Time
	Result     DriftResult
}

// Window is a sliding time window that retains recent drift results.
type Window struct {
	cfg     WindowConfig
	entries []WindowEntry
	now     func() time.Time
}

// NewWindow creates a Window with the given config.
func NewWindow(cfg WindowConfig) *Window {
	if cfg.Size <= 0 {
		cfg.Size = DefaultWindowConfig().Size
	}
	if cfg.MaxResults <= 0 {
		cfg.MaxResults = DefaultWindowConfig().MaxResults
	}
	return &Window{cfg: cfg, now: time.Now}
}

// Add inserts results into the window and evicts entries outside the window duration.
func (w *Window) Add(results []DriftResult) {
	now := w.now()
	for _, r := range results {
		w.entries = append(w.entries, WindowEntry{RecordedAt: now, Result: r})
	}
	w.evict(now)
	if len(w.entries) > w.cfg.MaxResults {
		w.entries = w.entries[len(w.entries)-w.cfg.MaxResults:]
	}
}

// Entries returns all results currently within the window.
func (w *Window) Entries() []WindowEntry {
	w.evict(w.now())
	out := make([]WindowEntry, len(w.entries))
	copy(out, w.entries)
	return out
}

// Results returns the DriftResult slice for all current window entries.
func (w *Window) Results() []DriftResult {
	entries := w.Entries()
	out := make([]DriftResult, len(entries))
	for i, e := range entries {
		out[i] = e.Result
	}
	return out
}

// Flush clears all entries from the window.
func (w *Window) Flush() {
	w.entries = nil
}

func (w *Window) evict(now time.Time) {
	cutoff := now.Add(-w.cfg.Size)
	start := 0
	for start < len(w.entries) && w.entries[start].RecordedAt.Before(cutoff) {
		start++
	}
	w.entries = w.entries[start:]
}
