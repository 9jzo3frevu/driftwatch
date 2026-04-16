package drift

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// HistoryEntry records a drift detection run.
type HistoryEntry struct {
	Timestamp time.Time    `json:"timestamp"`
	Results   []Result     `json:"results"`
	DriftCount int         `json:"drift_count"`
}

// History stores past drift detection results.
type History struct {
	mu      sync.RWMutex
	entries []HistoryEntry
	path    string
}

// NewHistory creates a History backed by the given file path.
func NewHistory(path string) *History {
	return &History{path: path}
}

// Record appends a new entry and persists to disk.
func (h *History) Record(results []Result) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	count := 0
	for _, r := range results {
		if r.Drift {
			count++
		}
	}

	entry := HistoryEntry{
		Timestamp:  time.Now().UTC(),
		Results:    results,
		DriftCount: count,
	}
	h.entries = append(h.entries, entry)
	return h.save()
}

// Entries returns a copy of all recorded entries.
func (h *History) Entries() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]HistoryEntry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Load reads persisted entries from disk.
func (h *History) Load() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := os.ReadFile(h.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.entries)
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0644)
}
