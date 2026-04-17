package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot captures the live state of a service at a point in time.
type Snapshot struct {
	ServiceID string            `json:"service_id"`
	CapturedAt time.Time        `json:"captured_at"`
	Values    map[string]string `json:"values"`
}

// NewSnapshot creates a Snapshot from a flat key-value map.
func NewSnapshot(serviceID string, values map[string]string) Snapshot {
	copy := make(map[string]string, len(values))
	for k, v := range values {
		copy[k] = v
	}
	return Snapshot{
		ServiceID:  serviceID,
		CapturedAt: time.Now().UTC(),
		Values:     copy,
	}
}

// SaveSnapshot writes a snapshot to dir/<serviceID>_snapshot.json.
func SaveSnapshot(dir string, s Snapshot) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir: %w", err)
	}
	path := filepath.Join(dir, s.ServiceID+"_snapshot.json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadSnapshot reads a previously saved snapshot for serviceID from dir.
func LoadSnapshot(dir, serviceID string) (Snapshot, error) {
	path := filepath.Join(dir, serviceID+"_snapshot.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: parse: %w", err)
	}
	return s, nil
}
