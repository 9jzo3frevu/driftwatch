package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWriter_WritesToStdout(t *testing.T) {
	w := NewWriter()
	if w.target != TargetStdout {
		t.Fatalf("expected TargetStdout, got %v", w.target)
	}
	if w.filePath != "" {
		t.Fatalf("expected empty filePath, got %q", w.filePath)
	}
}

func TestNewFileWriter_SetsPath(t *testing.T) {
	w := NewFileWriter("/tmp/report.txt")
	if w.target != TargetFile {
		t.Fatalf("expected TargetFile, got %v", w.target)
	}
	if w.filePath != "/tmp/report.txt" {
		t.Fatalf("expected /tmp/report.txt, got %q", w.filePath)
	}
}

func TestWriter_Open_Stdout(t *testing.T) {
	w := NewWriter()
	wc, err := w.Open()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer wc.Close()
	if wc == nil {
		t.Fatal("expected non-nil WriteCloser")
	}
	// Closing stdout wrapper should not error
	if err := wc.Close(); err != nil {
		t.Fatalf("Close() error: %v", err)
	}
}

func TestWriter_Open_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")

	w := NewFileWriter(path)
	wc, err := w.Open()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, writeErr := wc.Write([]byte("hello drift"))
	if writeErr != nil {
		t.Fatalf("write error: %v", writeErr)
	}
	wc.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(data), "hello drift") {
		t.Fatalf("expected file to contain 'hello drift', got %q", string(data))
	}
}

func TestWriter_Open_InvalidPath(t *testing.T) {
	w := NewFileWriter("/nonexistent/dir/report.txt")
	_, err := w.Open()
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
