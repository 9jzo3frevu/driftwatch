package source

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestPoller(t *testing.T, handler http.HandlerFunc, ttl time.Duration) (*Poller, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	fetcher := NewHTTPFetcher(2 * time.Second)
	cache := NewCache(ttl)
	poller := NewPoller(fetcher, cache, srv.URL, 50*time.Millisecond)
	return poller, srv
}

func TestPoller_Once_Success(t *testing.T) {
	payload := map[string]interface{}{"db.host": "localhost", "db.port": "5432"}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}
	poller, srv := newTestPoller(t, handler, time.Minute)
	defer srv.Close()

	result := poller.Once(context.Background())
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if result.Values["db.host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %q", result.Values["db.host"])
	}
}

func TestPoller_Once_PopulatesCache(t *testing.T) {
	payload := map[string]interface{}{"app.env": "production"}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}
	poller, srv := newTestPoller(t, handler, time.Minute)
	defer srv.Close()

	poller.Once(context.Background())
	v, ok := poller.cache.Get("app.env")
	if !ok || v != "production" {
		t.Errorf("expected cache hit for app.env=production, got ok=%v v=%q", ok, v)
	}
}

func TestPoller_CachedValues_Hit(t *testing.T) {
	poller := NewPoller(nil, NewCache(time.Minute), "", time.Second)
	poller.cache.Set("x", "1")
	poller.cache.Set("y", "2")

	out, err := poller.CachedValues([]string{"x", "y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["x"] != "1" || out["y"] != "2" {
		t.Errorf("unexpected values: %v", out)
	}
}

func TestPoller_CachedValues_Miss(t *testing.T) {
	poller := NewPoller(nil, NewCache(time.Minute), "", time.Second)
	_, err := poller.CachedValues([]string{"missing"})
	if err == nil {
		t.Error("expected error for cache miss")
	}
}

func TestPoller_Run_ReceivesResult(t *testing.T) {
	payload := map[string]interface{}{"k": "v"}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}
	poller, srv := newTestPoller(t, handler, time.Minute)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	ch := poller.Run(ctx)
	select {
	case result := <-ch:
		if result.Err != nil {
			t.Fatalf("unexpected error: %v", result.Err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for poll result")
	}
}
