package source

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader reads declared infrastructure-as-code state from files.
type Loader struct{}

// NewLoader creates a new Loader instance.
func NewLoader() *Loader {
	return &Loader{}
}

// LoadFile reads a JSON or YAML file and returns its contents as a flat
// map of string keys to string values.
func (l *Loader) LoadFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", path, err)
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		return parseJSON(data)
	case ".yaml", ".yml":
		return parseYAML(data)
	default:
		return nil, fmt.Errorf("unsupported file extension %q (want .json, .yaml, or .yml)", ext)
	}
}

func parseJSON(data []byte) (map[string]string, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return flattenMap("", raw), nil
}

func parseYAML(data []byte) (map[string]string, error) {
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}
	return flattenMap("", raw), nil
}

// flattenMap recursively flattens a nested map into dot-separated keys.
func flattenMap(prefix string, m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			for fk, fv := range flattenMap(key, val) {
				result[fk] = fv
			}
		default:
			result[key] = fmt.Sprintf("%v", val)
		}
	}
	return result
}
