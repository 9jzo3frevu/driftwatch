package source

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPFetcher_Fetch_OK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"env": "production", "replicas": "2", "db": {"host": "db.prod"}}`))
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher(5 * time.Second)
	got, err := fetcher.Fetch(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["env"] != "production" {
		t.Errorf("env: got %q, want %q", got["env"], "production")
	}
	if got["db.host"] != "db.prod" {
		t.Errorf("db.host: got %q, want %q", got["db.host"], "db.prod")
	}
}

func TestHTTPFetcher_Fetch_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher(5 * time.Second)
	_, err := fetcher.Fetch(server.URL)
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestHTTPFetcher_Fetch_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher(5 * time.Second)
	_, err := fetcher.Fetch(server.URL)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestHTTPFetcher_Fetch_DefaultTimeout(t *testing.T) {
	fetcher := NewHTTPFetcher(0)
	if fetcher.client.Timeout != 10*time.Second {
		t.Errorf("default timeout: got %v, want %v", fetcher.client.Timeout, 10*time.Second)
	}
}
