package config

import (
	"testing"
	"time"
)

func TestWebhookRaw_Build_Disabled(t *testing.T) {
	r := WebhookRaw{Enabled: false}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected disabled")
	}
}

func TestWebhookRaw_Build_MissingURL(t *testing.T) {
	r := WebhookRaw{Enabled: true}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for missing url")
	}
}

func TestWebhookRaw_Build_Valid(t *testing.T) {
	r := WebhookRaw{Enabled: true, URL: "https://example.com/hook", Secret: "s3cr3t", Timeout: "5s"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.URL != r.URL {
		t.Errorf("url mismatch: got %s", cfg.URL)
	}
	if cfg.Secret != r.Secret {
		t.Errorf("secret mismatch")
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("timeout mismatch: got %v", cfg.Timeout)
	}
}

func TestWebhookRaw_Build_DefaultTimeout(t *testing.T) {
	r := WebhookRaw{Enabled: true, URL: "https://example.com/hook"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Timeout != defaultWebhookTimeout {
		t.Errorf("expected default timeout, got %v", cfg.Timeout)
	}
}

func TestWebhookRaw_Build_InvalidTimeout(t *testing.T) {
	r := WebhookRaw{Enabled: true, URL: "https://example.com/hook", Timeout: "bad"}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for invalid timeout")
	}
}
