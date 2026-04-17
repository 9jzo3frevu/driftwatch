package alert

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/driftwatch/internal/config"
	"github.com/your-org/driftwatch/internal/drift"
)

func webhookResults() []drift.Result {
	return []drift.Result{
		{Key: "db.host", Declared: ptrStr("localhost"), Actual: ptrStr("remotehost")},
	}
}

func TestWebhookSender_Disabled(t *testing.T) {
	sender := NewWebhookSender(config.WebhookConfig{Enabled: false})
	if err := sender.Send(context.Background(), webhookResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWebhookSender_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := config.WebhookConfig{Enabled: true, URL: ts.URL, Timeout: 5 * time.Second}
	sender := NewWebhookSender(cfg)
	if err := sender.Send(context.Background(), webhookResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["drift_count"] == nil {
		t.Error("expected drift_count in payload")
	}
}

func TestWebhookSender_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cfg := config.WebhookConfig{Enabled: true, URL: ts.URL, Timeout: 5 * time.Second}
	sender := NewWebhookSender(cfg)
	if err := sender.Send(context.Background(), webhookResults()); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestWebhookSender_SignatureHeader(t *testing.T) {
	var sig string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sig = r.Header.Get("X-Driftwatch-Signature")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := config.WebhookConfig{Enabled: true, URL: ts.URL, Secret: "mysecret", Timeout: 5 * time.Second}
	sender := NewWebhookSender(cfg)
	if err := sender.Send(context.Background(), webhookResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sig == "" {
		t.Error("expected X-Driftwatch-Signature header")
	}
}
