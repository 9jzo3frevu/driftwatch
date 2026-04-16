package alert_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/alert"
)

func TestRawConfig_Build_Disabled(t *testing.T) {
	rc := alert.RawConfig{Enabled: false}
	cfg, err := rc.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.URL != "" {
		t.Errorf("expected empty URL for disabled alert, got %s", cfg.URL)
	}
}

func TestRawConfig_Build_MissingURL(t *testing.T) {
	rc := alert.RawConfig{Enabled: true}
	_, err := rc.Build()
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestRawConfig_Build_Valid(t *testing.T) {
	rc := alert.RawConfig{
		Enabled: true,
		Webhook: alert.RawWebhookConfig{
			URL:        "https://hooks.example.com/notify",
			TimeoutSec: 5,
			Headers:    map[string]string{"Authorization": "Bearer token"},
		},
	}
	cfg, err := rc.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.URL != "https://hooks.example.com/notify" {
		t.Errorf("unexpected URL: %s", cfg.URL)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", cfg.Timeout)
	}
	if cfg.Headers["Authorization"] != "Bearer token" {
		t.Errorf("expected auth header, got %v", cfg.Headers)
	}
}

func TestRawConfig_Build_DefaultTimeout(t *testing.T) {
	rc := alert.RawConfig{
		Enabled: true,
		Webhook: alert.RawWebhookConfig{
			URL: "https://hooks.example.com/notify",
		},
	}
	cfg, err := rc.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected default 10s timeout, got %v", cfg.Timeout)
	}
}
