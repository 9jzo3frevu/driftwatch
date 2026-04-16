package source

import (
	"testing"
	"time"
)

func TestCache_SetAndGet_Hit(t *testing.T) {
	c := NewCache(5 * time.Minute)
	data := map[string]string{"key": "value"}
	c.Set("svc", data)

	got, ok := c.Get("svc")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["key"] != "value" {
		t.Errorf("expected value, got %q", got["key"])
	}
}

func TestCache_Get_Miss(t *testing.T) {
	c := NewCache(5 * time.Minute)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestCache_Get_Expired(t *testing.T) {
	c := NewCache(1 * time.Millisecond)
	c.Set("svc", map[string]string{"a": "b"})
	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("svc")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := NewCache(5 * time.Minute)
	c.Set("svc", map[string]string{"x": "y"})
	c.Invalidate("svc")

	_, ok := c.Get("svc")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestCache_Flush(t *testing.T) {
	c := NewCache(5 * time.Minute)
	c.Set("a", map[string]string{})
	c.Set("b", map[string]string{})
	c.Flush()

	if c.Len() != 0 {
		t.Errorf("expected empty cache after flush, got %d entries", c.Len())
	}
}

func TestCache_Len(t *testing.T) {
	c := NewCache(5 * time.Minute)
	if c.Len() != 0 {
		t.Fatal("expected empty cache")
	}
	c.Set("svc1", map[string]string{})
	c.Set("svc2", map[string]string{})
	if c.Len() != 2 {
		t.Errorf("expected 2, got %d", c.Len())
	}
}
