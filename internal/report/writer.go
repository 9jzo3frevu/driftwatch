package report

import (
	"fmt"
	"io"
	"os"
)

// OutputTarget represents where report output should be written.
type OutputTarget int

const (
	TargetStdout OutputTarget = iota
	TargetFile
)

// Writer manages output destination for drift reports.
type Writer struct {
	target OutputTarget
	filePath string
}

// NewWriter creates a Writer that outputs to stdout.
func NewWriter() *Writer {
	return &Writer{target: TargetStdout}
}

// NewFileWriter creates a Writer that outputs to the given file path.
func NewFileWriter(path string) *Writer {
	return &Writer{target: TargetFile, filePath: path}
}

// Open returns a WriteCloser for the configured target.
// The caller is responsible for closing the returned WriteCloser.
func (w *Writer) Open() (io.WriteCloser, error) {
	switch w.target {
	case TargetFile:
		f, err := os.Create(w.filePath)
		if err != nil {
			return nil, fmt.Errorf("report writer: failed to open file %q: %w", w.filePath, err)
		}
		return f, nil
	default:
		return nopCloser{Writer: os.Stdout}, nil
	}
}

// nopCloser wraps an io.Writer with a no-op Close method.
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
