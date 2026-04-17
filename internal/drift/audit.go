package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditEntry records a single drift check event.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	DriftCount int      `json:"drift_count"`
	Score     int       `json:"score"`
	TriggeredBy string  `json:"triggered_by"`
}

// AuditLog appends and reads audit entries from a JSONL file.
type AuditLog struct {
	path string
}

// NewAuditLog creates an AuditLog writing to path.
func NewAuditLog(path string) *AuditLog {
	return &AuditLog{path: path}
}

// Record appends an entry to the audit log.
func (a *AuditLog) Record(entry AuditEntry) error {
	if err := os.MkdirAll(filepath.Dir(a.path), 0o755); err != nil {
		return fmt.Errorf("audit: mkdir: %w", err)
	}
	f, err := os.OpenFile(a.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode: %w", err)
	}
	return nil
}

// Entries reads all audit entries from the log file.
func (a *AuditLog) Entries() ([]AuditEntry, error) {
	f, err := os.Open(a.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: open: %w", err)
	}
	defer f.Close()
	var entries []AuditEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e AuditEntry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
