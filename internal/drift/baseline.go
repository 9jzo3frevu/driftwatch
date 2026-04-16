package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline holds a snapshot of expected key-value state.
type Baseline struct {
	CapturedAt time.Time         `json:"captured_at"`
	Values     map[string]string `json:"values"`
}

// NewBaseline creates a Baseline from the given flat key-value map.
func NewBaseline(values map[string]string) *Baseline {
	copy := make(map[string]string, len(values))
	for k, v := range values {
		copy[k] = v
	}
	return &Baseline{
		CapturedAt: time.Now().UTC(),
		Values:     copy,
	}
}

// SaveBaseline writes a Baseline to a JSON file at the given path.
func SaveBaseline(path string, b *Baseline) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal baseline: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write baseline: %w", err)
	}
	return nil
}

// LoadBaseline reads a Baseline from a JSON file at the given path.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline file not found: %s", path)
		}
		return nil, fmt.Errorf("read baseline: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("unmarshal baseline: %w", err)
	}
	return &b, nil
}
