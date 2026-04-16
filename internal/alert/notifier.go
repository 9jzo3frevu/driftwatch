package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Severity represents the alert severity level.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// WebhookConfig holds configuration for a webhook notifier.
type WebhookConfig struct {
	URL     string
	Timeout time.Duration
	Headers map[string]string
}

// Payload is the JSON body sent to the webhook.
type Payload struct {
	Service   string         `json:"service"`
	Severity  Severity       `json:"severity"`
	DriftCount int           `json:"drift_count"`
	Items     []drift.Result `json:"items"`
}

// Notifier sends drift alerts to a configured webhook.
type Notifier struct {
	cfg    WebhookConfig
	client *http.Client
}

// NewNotifier creates a Notifier with the given config.
func NewNotifier(cfg WebhookConfig) *Notifier {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &Notifier{
		cfg:    cfg,
		client: &http.Client{Timeout: timeout},
	}
}

// Notify sends a webhook notification if there are drift results.
func (n *Notifier) Notify(service string, results []drift.Result) error {
	if len(results) == 0 {
		return nil
	}

	severity := SeverityWarning
	if len(results) > 5 {
		severity = SeverityCritical
	}

	payload := Payload{
		Service:    service,
		Severity:   severity,
		DriftCount: len(results),
		Items:      results,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("alert: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("alert: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range n.cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("alert: send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("alert: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
