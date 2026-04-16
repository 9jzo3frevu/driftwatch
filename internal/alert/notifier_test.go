package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftwatch/internal/alert"
	"github.com/driftwatch/internal/drift"
)

func sampleResults(n int) []drift.Result {
	results := make([]drift.Result, n)
	for i := 0; i < n; i++ {
		results[i] = drift.Result{
			Key:      "key",
			Expected: "a",
			Actual:   "b",
		}
	}
	return results
}

func TestNotifier_Notify_NoDrift(t *testing.T) {
	n := alert.NewNotifier(alert.WebhookConfig{URL: "http://unused"})
	if err := n.Notify("svc", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNotifier_Notify_SendsPayload(t *testing.T) {
	var received alert.Payload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := alert.NewNotifier(alert.WebhookConfig{URL: ts.URL})
	results := sampleResults(2)
	if err := n.Notify("my-service", results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Service != "my-service" {
		t.Errorf("expected service my-service, got %s", received.Service)
	}
	if received.DriftCount != 2 {
		t.Errorf("expected drift_count 2, got %d", received.DriftCount)
	}
	if received.Severity != alert.SeverityWarning {
		t.Errorf("expected warning severity, got %s", received.Severity)
	}
}

func TestNotifier_Notify_CriticalSeverity(t *testing.T) {
	var received alert.Payload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := alert.NewNotifier(alert.WebhookConfig{URL: ts.URL})
	if err := n.Notify("svc", sampleResults(6)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Severity != alert.SeverityCritical {
		t.Errorf("expected critical, got %s", received.Severity)
	}
}

func TestNotifier_Notify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := alert.NewNotifier(alert.WebhookConfig{URL: ts.URL})
	if err := n.Notify("svc", sampleResults(1)); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestNotifier_DefaultTimeout(t *testing.T) {
	n := alert.NewNotifier(alert.WebhookConfig{URL: "http://unused"})
	_ = n
	// Ensure construction with zero timeout doesn't panic
	n2 := alert.NewNotifier(alert.WebhookConfig{URL: "http://unused", Timeout: 5 * time.Second})
	_ = n2
}
