package drift

import (
	"testing"
	"time"
)

func throttleResults() []DriftResult {
	return []DriftResult{
		{Service: "api", Key: "db.host", Expected: "prod-db", Actual: ptrStr("staging-db")},
		{Service: "api", Key: "db.port", Expected: "5432", Actual: ptrStr("5433")},
		{Service: "worker", Key: "queue.url", Expected: "amqp://prod", Actual: ptrStr("amqp://dev")},
	}
}

func TestThrottler_FirstOccurrence_Passes(t *testing.T) {
	th := NewThrottler(ThrottleConfig{Window: time.Minute, MaxPerWindow: 2})
	results := throttleResults()
	out := th.Apply(results)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestThrottler_ExceedsMax_Suppressed(t *testing.T) {
	th := NewThrottler(ThrottleConfig{Window: time.Minute, MaxPerWindow: 2})
	single := throttleResults()[:1]

	// First two calls should pass
	out1 := th.Apply(single)
	out2 := th.Apply(single)
	// Third should be suppressed
	out3 := th.Apply(single)

	if len(out1) != 1 {
		t.Errorf("pass 1: expected 1, got %d", len(out1))
	}
	if len(out2) != 1 {
		t.Errorf("pass 2: expected 1, got %d", len(out2))
	}
	if len(out3) != 0 {
		t.Errorf("pass 3: expected 0 (suppressed), got %d", len(out3))
	}
}

func TestThrottler_WindowExpiry_Resets(t *testing.T) {
	th := NewThrottler(ThrottleConfig{Window: 10 * time.Millisecond, MaxPerWindow: 1})
	single := throttleResults()[:1]

	// First call passes, second is suppressed
	th.Apply(single)
	out := th.Apply(single)
	if len(out) != 0 {
		t.Fatalf("expected suppression before window expiry, got %d", len(out))
	}

	// Wait for window to expire
	time.Sleep(20 * time.Millisecond)

	out2 := th.Apply(single)
	if len(out2) != 1 {
		t.Fatalf("expected result after window expiry, got %d", len(out2))
	}
}

func TestThrottler_DifferentKeys_Independent(t *testing.T) {
	th := NewThrottler(ThrottleConfig{Window: time.Minute, MaxPerWindow: 1})
	results := throttleResults() // 3 distinct service/key combos

	// Each key gets its own bucket; all should pass on first call
	out := th.Apply(results)
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestThrottler_Reset_ClearsState(t *testing.T) {
	th := NewThrottler(ThrottleConfig{Window: time.Minute, MaxPerWindow: 1})
	single := throttleResults()[:1]

	th.Apply(single)
	th.Apply(single) // second is suppressed

	th.Reset()

	out := th.Apply(single)
	if len(out) != 1 {
		t.Fatalf("expected 1 after reset, got %d", len(out))
	}
}

func TestThrottler_DefaultsOnZeroConfig(t *testing.T) {
	th := NewThrottler(ThrottleConfig{})
	def := DefaultThrottleConfig()
	if th.cfg.Window != def.Window {
		t.Errorf("expected default window %v, got %v", def.Window, th.cfg.Window)
	}
	if th.cfg.MaxPerWindow != def.MaxPerWindow {
		t.Errorf("expected default max %d, got %d", def.MaxPerWindow, th.cfg.MaxPerWindow)
	}
}
