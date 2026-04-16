package alert

import (
	"fmt"
	"time"
)

// RawConfig is the alert section parsed from the top-level config file.
type RawConfig struct {
	Enabled bool              `yaml:"enabled"`
	Webhook RawWebhookConfig  `yaml:"webhook"`
}

// RawWebhookConfig holds raw YAML fields for the webhook.
type RawWebhookConfig struct {
	URL        string            `yaml:"url"`
	TimeoutSec int               `yaml:"timeout_seconds"`
	Headers    map[string]string `yaml:"headers"`
}

// Build validates and converts RawConfig into a WebhookConfig.
func (rc RawConfig) Build() (WebhookConfig, error) {
	if !rc.Enabled {
		return WebhookConfig{}, nil
	}
	if rc.Webhook.URL == "" {
		return WebhookConfig{}, fmt.Errorf("alert: webhook url is required when alerts are enabled")
	}
	timeout := time.Duration(rc.Webhook.TimeoutSec) * time.Second
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return WebhookConfig{
		URL:     rc.Webhook.URL,
		Timeout: timeout,
		Headers: rc.Webhook.Headers,
	}, nil
}
