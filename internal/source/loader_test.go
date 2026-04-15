package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoader_LoadFile_JSON(t *testing.T) {
	content := `{"service": "api", "replicas": "3", "db": {"host": "localhost", "port": "5432"}}`
	path := writeTempFile(t, "config.json", content)

	loader := NewLoader()
	got, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{
		"service":  "api",
		"replicas": "3",
		"db.host":  "localhost",
		"db.port":  "5432",
	}
	assertMapsEqual(t, expected, got)
}

func TestLoader_LoadFile_YAML(t *testing.T) {
	content := "service: api\nreplicas: 3\ndb:\n  host: localhost\n  port: 5432\n"
	path := writeTempFile(t, "config.yaml", content)

	loader := NewLoader()
	got, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["service"] != "api" {
		t.Errorf("service: got %q, want %q", got["service"], "api")
	}
	if got["db.host"] != "localhost" {
		t.Errorf("db.host: got %q, want %q", got["db.host"], "localhost")
	}
}

func TestLoader_LoadFile_UnsupportedExtension(t *testing.T) {
	path := writeTempFile(t, "config.toml", "key = \"value\"")

	loader := NewLoader()
	_, err := loader.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for unsupported extension, got nil")
	}
}

func TestLoader_LoadFile_NotFound(t *testing.T) {
	loader := NewLoader()
	_, err := loader.LoadFile("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func assertMapsEqual(t *testing.T, expected, got map[string]string) {
	t.Helper()
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("key %q: got %q, want %q", k, got[k], v)
		}
	}
	for k := range got {
		if _, ok := expected[k]; !ok {
			t.Errorf("unexpected key %q in result", k)
		}
	}
}
