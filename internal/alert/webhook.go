package alert

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/driftwatch/internal/config"
	"github.com/your-org/driftwatch/internal/drift"
)

// WebhookSender sends drift results to a webhook endpoint.
type WebhookSender struct {
	cfg    config.WebhookConfig
	client *http.Client
}

// NewWebhookSender creates a WebhookSender from the given config.
func NewWebhookSender(cfg config.WebhookConfig) *WebhookSender {
	return &WebhookSender{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

type webhookPayload struct {
	Timestamp string         `json:"timestamp"`
	DriftCount int           `json:"drift_count"`
	Results   []drift.Result `json:"results"`
}

// Send posts drift results to the configured webhook URL.
func (w *WebhookSender) Send(ctx context.Context, results []drift.Result) error {
	if !w.cfg.Enabled || len(results) == 0 {
		return nil
	}
	payload := webhookPayload{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		DriftCount: len(results),
		Results:    results,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if w.cfg.Secret != "" {
		sig := sign(body, w.cfg.Secret)
		req.Header.Set("X-Driftwatch-Signature", sig)
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func sign(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
