package report

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Format defines the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes drift results to an output stream.
type Formatter struct {
	Format Format
	Writer io.Writer
}

// NewFormatter creates a Formatter with the given format and writer.
func NewFormatter(format Format, w io.Writer) *Formatter {
	return &Formatter{Format: format, Writer: w}
}

// Write outputs the drift results according to the configured format.
func (f *Formatter) Write(results []drift.DriftResult) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(results)
	default:
		return f.writeText(results)
	}
}

func (f *Formatter) writeText(results []drift.DriftResult) error {
	if len(results) == 0 {
		fmt.Fprintln(f.Writer, "✓ No drift detected.")
		return nil
	}

	w := tabwriter.NewWriter(f.Writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tTYPE\tEXPECTED\tACTUAL")
	fmt.Fprintln(w, strings.Repeat("-", 60))
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			r.Key,
			r.ChangeType,
			valOrEmpty(r.Expected),
			valOrEmpty(r.Actual),
		)
	}
	return w.Flush()
}

func (f *Formatter) writeJSON(results []drift.DriftResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(f.Writer, `{"drift": []}`)
		return err
	}

	fmt.Fprintln(f.Writer, `{"drift": [`)
	for i, r := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(f.Writer, "  {\"key\": %q, \"type\": %q, \"expected\": %q, \"actual\": %q}%s\n",
			r.Key, r.ChangeType, valOrEmpty(r.Expected), valOrEmpty(r.Actual), comma)
	}
	_, err := fmt.Fprintln(f.Writer, `]}`)
	return err
}

func valOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
