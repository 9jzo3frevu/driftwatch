package drift

import (
	"context"
	"errors"
	"testing"
)

func pipelineResults() []DriftResult {
	return []DriftResult{
		{Key: "app.timeout", Declared: "30s", Live: ptrStr("60s"), Drifted: true},
		{Key: "app.replicas", Declared: "3", Live: ptrStr("3"), Drifted: false},
	}
}

func TestPipeline_Empty_ReturnsInput(t *testing.T) {
	p := NewPipeline()
	results := pipelineResults()
	out, err := p.Run(context.Background(), results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(out))
	}
}

func TestPipeline_StagesAppliedInOrder(t *testing.T) {
	var order []string
	makeStage := func(name string) Stage {
		return func(r []DriftResult) ([]DriftResult, error) {
			order = append(order, name)
			return r, nil
		}
	}
	p := NewPipeline().
		Register("first", makeStage("first")).
		Register("second", makeStage("second")).
		Register("third", makeStage("third"))

	_, err := p.Run(context.Background(), pipelineResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 3 || order[0] != "first" || order[1] != "second" || order[2] != "third" {
		t.Errorf("unexpected stage order: %v", order)
	}
}

func TestPipeline_StageError_ShortCircuits(t *testing.T) {
	ran := false
	p := NewPipeline().
		Register("fail", func(r []DriftResult) ([]DriftResult, error) {
			return nil, errors.New("stage failure")
		}).
		Register("after", func(r []DriftResult) ([]DriftResult, error) {
			ran = true
			return r, nil
		})

	_, err := p.Run(context.Background(), pipelineResults())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if ran {
		t.Error("stage after failure should not have run")
	}
}

func TestPipeline_ContextCancelled_StopsEarly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := NewPipeline().
		Register("any", func(r []DriftResult) ([]DriftResult, error) {
			return r, nil
		})

	_, err := p.Run(ctx, pipelineResults())
	if err == nil {
		t.Fatal("expected cancellation error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}

func TestPipeline_StageCount(t *testing.T) {
	p := NewPipeline().
		Register("a", func(r []DriftResult) ([]DriftResult, error) { return r, nil }).
		Register("b", func(r []DriftResult) ([]DriftResult, error) { return r, nil })
	if p.StageCount() != 2 {
		t.Errorf("expected 2 stages, got %d", p.StageCount())
	}
}
