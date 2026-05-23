package drift

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Checkpoint records the last successful drift run for a given service.
type Checkpoint struct {
	Service   string    `json:"service"`
	RunAt     time.Time `json:"run_at"`
	DriftCount int      `json:"drift_count"`
}

// CheckpointStore persists and loads Checkpoint records to disk.
type CheckpointStore struct {
	dir string
}

// NewCheckpointStore creates a CheckpointStore that writes to dir.
func NewCheckpointStore(dir string) *CheckpointStore {
	return &CheckpointStore{dir: dir}
}

// Save writes the checkpoint for the given service to disk.
func (s *CheckpointStore) Save(cp Checkpoint) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath(cp.Service), data, 0o644)
}

// Load returns the checkpoint for the given service.
// Returns a zero Checkpoint and no error when not found.
func (s *CheckpointStore) Load(service string) (Checkpoint, error) {
	data, err := os.ReadFile(s.filePath(service))
	if errors.Is(err, os.ErrNotExist) {
		return Checkpoint{}, nil
	}
	if err != nil {
		return Checkpoint{}, err
	}
	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return Checkpoint{}, err
	}
	return cp, nil
}

// Delete removes the checkpoint file for the given service.
func (s *CheckpointStore) Delete(service string) error {
	err := os.Remove(s.filePath(service))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *CheckpointStore) filePath(service string) string {
	return filepath.Join(s.dir, service+".checkpoint.json")
}
