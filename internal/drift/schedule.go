package drift

import (
	"context"
	"log"
	"time"
)

// ScheduleConfig holds configuration for scheduled drift checks.
type ScheduleConfig struct {
	Interval time.Duration
	MaxRuns  int // 0 means unlimited
}

// Scheduler runs drift detection on a fixed interval.
type Scheduler struct {
	cfg     ScheduleConfig
	detect  func(ctx context.Context) ([]Result, error)
	onDrift func(results []Result)
}

// NewScheduler creates a Scheduler with the given config and callbacks.
func NewScheduler(cfg ScheduleConfig, detect func(ctx context.Context) ([]Result, error), onDrift func([]Result)) *Scheduler {
	return &Scheduler{cfg: cfg, detect: detect, onDrift: onDrift}
}

// Run starts the scheduler loop, blocking until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()
	runs := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			results, err := s.detect(ctx)
			if err != nil {
				log.Printf("scheduler: detect error: %v", err)
			} else if len(results) > 0 && s.onDrift != nil {
				s.onDrift(results)
			}
			runs++
			if s.cfg.MaxRuns > 0 && runs >= s.cfg.MaxRuns {
				return
			}
		}
	}
}
