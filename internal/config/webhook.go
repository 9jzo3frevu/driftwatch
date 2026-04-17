package config

import (
	"errors"
	"fmt"
	"time"
)

// WebhookRaw holds raw webhook configuration from YAML.
type WebhookRaw struct {
	Enabled bool   `yaml:"enabled"`
	URL     string `yaml:"url"`
	Secret  string `yaml:"secret"`
	Timeout string `yaml:"timeout"`
}

// WebhookConfig is the validated webhook configuration.
type WebhookConfig struct {
	Enabled bool
	URL     string
	Secret  string
	Timeout time.Duration
}

const defaultWebhookTimeout = 10 * time.Second

// Build validates and builds a WebhookConfig from raw values.
func (r WebhookRaw) Build() (WebhookConfig, error) {
	if !r.Enabled {
		return WebhookConfig{Enabled: false}, nil
	}
	if r.URL == "" {
		return WebhookConfig{}, errors.New("webhook: url is required when enabled")
	}
	timeout := defaultWebhookTimeout
	if r.Timeout != "" {
		d, err := time.ParseDuration(r.Timeout)
		if err != nil {
			return WebhookConfig{}, fmt.Errorf("webhook: invalid timeout %q: %w", r.Timeout, err)
		}
		timeout = d
	}
	return WebhookConfig{
		Enabled: true,
		URL:     r.URL,
		Secret:  r.Secret,
		Timeout: timeout,
	}, nil
}
