package source

import (
	"context"
	"fmt"
	"time"
)

// PollResult holds the result of a single poll attempt.
type PollResult struct {
	Values map[string]string
	FetchedAt time.Time
	Err error
}

// Poller periodically fetches live configuration from a remote endpoint.
type Poller struct {
	fetcher *HTTPFetcher
	cache *Cache
	interval time.Duration
	endpoint string
}

// NewPoller creates a Poller with the given fetcher, cache, endpoint, and poll interval.
func NewPoller(fetcher *HTTPFetcher, cache *Cache, endpoint string, interval time.Duration) *Poller {
	return &Poller{
		fetcher: fetcher,
		cache: cache,
		interval: interval,
		endpoint: endpoint,
	}
}

// Once performs a single fetch, populates the cache, and returns the result.
func (p *Poller) Once(ctx context.Context) PollResult {
	values, err := p.fetcher.Fetch(ctx, p.endpoint)
	result := PollResult{
		FetchedAt: time.Now(),
		Err: err,
	}
	if err == nil {
		result.Values = values
		for k, v := range values {
			p.cache.Set(k, v)
		}
	}
	return result
}

// Run polls on the given interval until ctx is cancelled, sending results to the returned channel.
func (p *Poller) Run(ctx context.Context) <-chan PollResult {
	ch := make(chan PollResult, 1)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(p.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				result := p.Once(ctx)
			ect {
		default:
			 if consumer is slow
				}
		}
	}()
	return ch
}

// CachedValues returns all currently cached values or an error if the cache is empty.
func (p *Poller) CachedValues(keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		v, ok := p.cache.Get(k)
		if !ok {
			return nil, fmt.Errorf("cache miss for key %q", k)
		}
		out[k] = v
	}
	return out, nil
}
