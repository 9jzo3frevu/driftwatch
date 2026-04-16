package source

import (
	"sync"
	"time"
)

// CacheEntry holds a cached flat map result with an expiry timestamp.
type CacheEntry struct {
	Data      map[string]string
	FetchedAt time.Time
	TTL       time.Duration
}

// IsExpired returns true if the cache entry has exceeded its TTL.
func (e *CacheEntry) IsExpired() bool {
	return time.Since(e.FetchedAt) > e.TTL
}

// Cache is a simple in-memory, thread-safe store for fetched source data.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	defaultTTL time.Duration
}

// NewCache creates a Cache with the given default TTL.
func NewCache(defaultTTL time.Duration) *Cache {
	return &Cache{
		entries:    make(map[string]*CacheEntry),
		defaultTTL: defaultTTL,
	}
}

// Get returns the cached data for key if present and not expired.
// The second return value indicates a cache hit.
func (c *Cache) Get(key string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok || entry.IsExpired() {
		return nil, false
	}
	return entry.Data, true
}

// Set stores data under key using the cache's default TTL.
func (c *Cache) Set(key string, data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = &CacheEntry{
		Data:      data,
		FetchedAt: time.Now(),
		TTL:       c.defaultTTL,
	}
}

// Invalidate removes a single entry from the cache.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}
