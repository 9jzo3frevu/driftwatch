package drift

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestScheduler_RunsDetect(t *testing.T) {
	var count int32
	detect := func(_ context.Context) ([]Result, error) {
		atomic.AddInt32(&count, 1)
		return nil, nil
	}
	cfg := ScheduleConfig{Interval: 20 * time.Millisecond, MaxRuns: 3}
	s := NewScheduler(cfg, detect, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	s.Run(ctx)
	if atomic.LoadInt32(&count) != 3 {
		t.Fatalf("expected 3 runs, got %d", count)
	}
}

func TestScheduler_CallsOnDrift(t *testing.T) {
	drifted := []Result{{Key: "x", Declared: ptrStr("a"), Live: ptrStr("b")}}
	var called int32
	detect := func(_ context.Context) ([]Result, error) { return drifted, nil }
	onDrift := func(r []Result) { atomic.AddInt32(&called, 1) }
	cfg := ScheduleConfig{Interval: 20 * time.Millisecond, MaxRuns: 2}
	s := NewScheduler(cfg, detect, onDrift)
	s.Run(context.Background())
	if atomic.LoadInt32(&called) != 2 {
		t.Fatalf("expected onDrift called 2 times, got %d", called)
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	var count int32
	detect := func(_ context.Context) ([]Result, error) {
		atomic.AddInt32(&count, 1)
		return nil, nil
	}
	cfg := ScheduleConfig{Interval: 10 * time.Millisecond}
	s := NewScheduler(cfg, detect, nil)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(35 * time.Millisecond)
		cancel()
	}()
	s.Run(ctx)
	if atomic.LoadInt32(&count) == 0 {
		t.Fatal("expected at least one run before cancel")
	}
}

func TestScheduler_NoDriftNoCallback(t *testing.T) {
	detect := func(_ context.Context) ([]Result, error) { return nil, nil }
	called := false
	onDrift := func(_ []Result) { called = true }
	cfg := ScheduleConfig{Interval: 20 * time.Millisecond, MaxRuns: 2}
	s := NewScheduler(cfg, detect, onDrift)
	s.Run(context.Background())
	if called {
		t.Fatal("onDrift should not be called when no drift")
	}
}
