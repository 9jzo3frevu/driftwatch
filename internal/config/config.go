package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level driftwatch configuration.
type Config struct {
	Sources  SourcesConfig  `yaml:"sources"`
	Report   ReportConfig   `yaml:"report"`
	Detector DetectorConfig `yaml:"detector"`
}

// SourcesConfig defines where declared and live state come from.
type SourcesConfig struct {
	DeclaredFile string        `yaml:"declared_file"`
	LiveURL      string        `yaml:"live_url"`
	Timeout      time.Duration `yaml:"timeout"`
}

// ReportConfig controls output format and destination.
type ReportConfig struct {
	Format string `yaml:"format"` // "text" or "json"
	Output string `yaml:"output"` // file path or "stdout"
}

// DetectorConfig holds drift detection options.
type DetectorConfig struct {
	IgnoreKeys []string `yaml:"ignore_keys"`
}

// Load reads a YAML config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	if cfg.Sources.Timeout == 0 {
		cfg.Sources.Timeout = 10 * time.Second
	}
	if cfg.Report.Format == "" {
		cfg.Report.Format = "text"
	}
	if cfg.Report.Output == "" {
		cfg.Report.Output = "stdout"
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Sources.DeclaredFile == "" {
		return fmt.Errorf("sources.declared_file is required")
	}
	if c.Sources.LiveURL == "" {
		return fmt.Errorf("sources.live_url is required")
	}
	switch c.Report.Format {
	case "", "text", "json":
		// valid
	default:
		return fmt.Errorf("report.format must be \"text\" or \"json\", got %q", c.Report.Format)
	}
	return nil
}
