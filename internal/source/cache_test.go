package source

import (
	"testing"
	"time"
)

func TestCache_SetAndGet_Hit(t *testing.T) {
	c := NewCache(5 * time.Minute)
	data := map[string]string{"key": "value"}
	c.Set("svc-a", data)

	got, ok := c.Get("svc-a")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if got["key"] != "value" {
		t.Errorf("expected 'value', got %q", got["key"])
	}
}

func TestCache_Get_Miss(t *testing.T) {
	c := NewCache(5 * time.Minute)

	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestCache_Get_Expired(t *testing.T) {
	c := NewCache(1 * time.Millisecond)
	c.Set("svc-b", map[string]string{"x": "1"})

	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("svc-b")
	if ok {
		t.Fatal("expected expired cache miss, got hit")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := NewCache(5 * time.Minute)
	c.Set("svc-c", map[string]string{"a": "b"})
	c.Invalidate("svc-c")

	_, ok := c.Get("svc-c")
	if ok {
		t.Fatal("expected miss after invalidation")
	}
}

func TestCache_Flush(t *testing.T) {
	c := NewCache(5 * time.Minute)
	c.Set("svc-d", map[string]string{"p": "q"})
	c.Set("svc-e", map[string]string{"r": "s"})
	c.Flush()

	for _, key := range []string{"svc-d", "svc-e"} {
		if _, ok := c.Get(key); ok {
			t.Errorf("expected miss for %q after flush", key)
		}
	}
}

func TestCacheEntry_IsExpired(t *testing.T) {
	entry := &CacheEntry{
		FetchedAt: time.Now().Add(-10 * time.Second),
		TTL:       5 * time.Second,
	}
	if !entry.IsExpired() {
		t.Error("expected entry to be expired")
	}

	fresh := &CacheEntry{
		FetchedAt: time.Now(),
		TTL:       5 * time.Minute,
	}
	if fresh.IsExpired() {
		t.Error("expected fresh entry to not be expired")
	}
}
