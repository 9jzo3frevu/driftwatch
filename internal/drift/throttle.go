package drift

import (
	"sync"
	"time"
)

// ThrottleConfig controls how alert throttling behaves.
type ThrottleConfig struct {
	// Window is the duration during which repeated drift on the same key is suppressed.
	Window time.Duration
	// MaxPerWindow is the maximum number of alerts allowed per key per window.
	MaxPerWindow int
}

// DefaultThrottleConfig returns sensible defaults for throttling.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		Window:       5 * time.Minute,
		MaxPerWindow: 3,
	}
}

type throttleEntry struct {
	count     int
	windowEnd time.Time
}

// Throttler suppresses drift results that have been seen too frequently
// within a rolling time window.
type Throttler struct {
	cfg     ThrottleConfig
	mu      sync.Mutex
	buckets map[string]*throttleEntry
}

// NewThrottler creates a Throttler with the given config.
func NewThrottler(cfg ThrottleConfig) *Throttler {
	if cfg.Window <= 0 {
		cfg.Window = DefaultThrottleConfig().Window
	}
	if cfg.MaxPerWindow <= 0 {
		cfg.MaxPerWindow = DefaultThrottleConfig().MaxPerWindow
	}
	return &Throttler{
		cfg:     cfg,
		buckets: make(map[string]*throttleEntry),
	}
}

// Apply filters out drift results that exceed the per-key rate limit.
// Results that are within limits are returned; suppressed ones are dropped.
func (t *Throttler) Apply(results []DriftResult) []DriftResult {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	var out []DriftResult
	for _, r := range results {
		key := r.Service + "/" + r.Key
		ent, ok := t.buckets[key]
		if !ok || now.After(ent.windowEnd) {
			t.buckets[key] = &throttleEntry{
				count:     1,
				windowEnd: now.Add(t.cfg.Window),
			}
			out = append(out, r)
			continue
		}
		if ent.count < t.cfg.MaxPerWindow {
			ent.count++
			out = append(out, r)
		}
		// else: suppressed — too many within window
	}
	return out
}

// Reset clears all throttle state, useful for testing or config reloads.
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buckets = make(map[string]*throttleEntry)
}
