package drift

import (
	"context"
	"fmt"
)

// Stage is a function that transforms a slice of DriftResults.
type Stage func([]DriftResult) ([]DriftResult, error)

// Pipeline chains multiple processing stages applied sequentially to drift results.
type Pipeline struct {
	stages []namedStage
}

type namedStage struct {
	name  string
	apply Stage
}

// NewPipeline returns an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Register adds a named stage to the pipeline.
func (p *Pipeline) Register(name string, s Stage) *Pipeline {
	p.stages = append(p.stages, namedStage{name: name, apply: s})
	return p
}

// Run executes all stages in order, short-circuiting on error.
// The context is checked between stages for cancellation.
func (p *Pipeline) Run(ctx context.Context, results []DriftResult) ([]DriftResult, error) {
	current := results
	for _, s := range p.stages {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("pipeline cancelled before stage %q: %w", s.name, ctx.Err())
		default:
		}
		var err error
		current, err = s.apply(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %q: %w", s.name, err)
		}
	}
	return current, nil
}

// StageCount returns the number of registered stages.
func (p *Pipeline) StageCount() int {
	return len(p.stages)
}
