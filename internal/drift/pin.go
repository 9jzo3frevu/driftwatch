package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PinnedKey represents a drift result key that has been pinned to its current
// live value, suppressing future drift alerts until the pin expires or is removed.
type PinnedKey struct {
	Key       string    `json:"key"`
	Service   string    `json:"service"`
	PinnedAt  time.Time `json:"pinned_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Reason    string    `json:"reason,omitempty"`
}

// IsExpired reports whether the pin has passed its expiry time.
func (p PinnedKey) IsExpired(now time.Time) bool {
	if p.ExpiresAt.IsZero() {
		return false
	}
	return now.After(p.ExpiresAt)
}

// PinStore persists and retrieves pinned keys.
type PinStore struct {
	path string
}

// NewPinStore creates a PinStore backed by the given file path.
func NewPinStore(path string) *PinStore {
	return &PinStore{path: path}
}

// Load reads all pins from disk, discarding expired entries.
func (s *PinStore) Load(now time.Time) ([]PinnedKey, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("pin store read: %w", err)
	}
	var pins []PinnedKey
	if err := json.Unmarshal(data, &pins); err != nil {
		return nil, fmt.Errorf("pin store parse: %w", err)
	}
	active := pins[:0]
	for _, p := range pins {
		if !p.IsExpired(now) {
			active = append(active, p)
		}
	}
	return active, nil
}

// Save writes the given pins to disk, creating parent directories as needed.
func (s *PinStore) Save(pins []PinnedKey) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("pin store mkdir: %w", err)
	}
	data, err := json.MarshalIndent(pins, "", "  ")
	if err != nil {
		return fmt.Errorf("pin store encode: %w", err)
	}
	return os.WriteFile(s.path, data, 0o644)
}

// ApplyPins removes any DriftResult whose (Service, Key) pair is currently pinned.
func ApplyPins(results []DriftResult, pins []PinnedKey) []DriftResult {
	index := make(map[string]struct{}, len(pins))
	for _, p := range pins {
		index[p.Service+"\x00"+p.Key] = struct{}{}
	}
	out := results[:0:0]
	for _, r := range results {
		if _, pinned := index[r.Service+"\x00"+r.Key]; !pinned {
			out = append(out, r)
		}
	}
	return out
}
