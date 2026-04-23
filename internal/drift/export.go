package drift

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportFormat defines the output format for exported drift results.
type ExportFormat string

const (
	ExportCSV  ExportFormat = "csv"
	ExportJSON ExportFormat = "json"
)

// ExportConfig controls how results are exported.
type ExportConfig struct {
	Format    ExportFormat
	Timestamp bool
}

// Exporter writes drift results to an io.Writer in the configured format.
type Exporter struct {
	cfg ExportConfig
}

// NewExporter creates an Exporter with the given config.
func NewExporter(cfg ExportConfig) *Exporter {
	if cfg.Format == "" {
		cfg.Format = ExportCSV
	}
	return &Exporter{cfg: cfg}
}

// Export writes results to w in the configured format.
func (e *Exporter) Export(w io.Writer, results []DriftResult) error {
	switch e.cfg.Format {
	case ExportJSON:
		return e.writeJSON(w, results)
	case ExportCSV:
		return e.writeCSV(w, results)
	default:
		return fmt.Errorf("unsupported export format: %s", e.cfg.Format)
	}
}

type exportRow struct {
	Key      string  `json:"key"`
	Service  string  `json:"service"`
	Declared string  `json:"declared"`
	Live     string  `json:"live"`
	Severity string  `json:"severity"`
	At       *string `json:"exported_at,omitempty"`
}

func (e *Exporter) buildRows(results []DriftResult) []exportRow {
	var ts *string
	if e.cfg.Timestamp {
		s := time.Now().UTC().Format(time.RFC3339)
		ts = &s
	}
	rows := make([]exportRow, 0, len(results))
	for _, r := range results {
		rows = append(rows, exportRow{
			Key:      r.Key,
			Service:  r.Service,
			Declared: ptrStr(r.Declared),
			Live:     ptrStr(r.Live),
			Severity: string(r.Severity),
			At:       ts,
		})
	}
	return rows
}

func (e *Exporter) writeJSON(w io.Writer, results []DriftResult) error {
	rows := e.buildRows(results)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

func (e *Exporter) writeCSV(w io.Writer, results []DriftResult) error {
	cw := csv.NewWriter(w)
	header := []string{"key", "service", "declared", "live", "severity"}
	if e.cfg.Timestamp {
		header = append(header, "exported_at")
	}
	if err := cw.Write(header); err != nil {
		return err
	}
	for _, row := range e.buildRows(results) {
		rec := []string{row.Key, row.Service, row.Declared, row.Live, row.Severity}
		if e.cfg.Timestamp {
			rec = append(rec, ptrStr(row.At))
		}
		if err := cw.Write(rec); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
