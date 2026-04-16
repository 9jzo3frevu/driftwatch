package source

import (
	"sync"
	"time"
)

type cacheEntry struct {
	data      map[string]string
	fetchedAt time.Time
	ttl       time.Duration
}

func (e *cacheEntry) expired() bool {
	return time.Since(e.fetchedAt) > e.ttl
}

// Cache stores fetched remote configs with TTL-based expiry.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*cacheEntry
	defaultTTL time.Duration
}

// NewCache creates a Cache with the given default TTL.
func NewCache(defaultTTL time.Duration) *Cache {
	return &Cache{
		entries:    make(map[string]*cacheEntry),
		defaultTTL: defaultTTL,
	}
}

// Set stores data for key with the default TTL.
func (c *Cache) Set(key string, data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &cacheEntry{
		data:      data,
		fetchedAt: time.Now(),
		ttl:       c.defaultTTL,
	}
}

// Get retrieves data for key. Returns nil, false if missing or expired.
func (c *Cache) Get(key string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || e.expired() {
		return nil, false
	}
	return e.data, true
}

// Invalidate removes a single key from the cache.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*cacheEntry)
}

// Len returns the number of entries currently in the cache (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
