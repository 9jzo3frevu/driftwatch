package config

import (
	"testing"
)

func boolPtrP(b bool) *bool { return &b }

func TestPipelineRaw_Build_Empty(t *testing.T) {
	r := &PipelineRaw{}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.EnabledStages) != 0 {
		t.Errorf("expected no stages, got %v", cfg.EnabledStages)
	}
}

func TestPipelineRaw_Build_NilIsEmpty(t *testing.T) {
	var r *PipelineRaw
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil || len(cfg.EnabledStages) != 0 {
		t.Error("expected empty config from nil raw")
	}
}

func TestPipelineRaw_Build_MissingName(t *testing.T) {
	r := &PipelineRaw{
		Stages: []PipelineStageRaw{{Name: ""}},
	}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestPipelineRaw_Build_DuplicateName(t *testing.T) {
	r := &PipelineRaw{
		Stages: []PipelineStageRaw{
			{Name: "filter"},
			{Name: "filter"},
		},
	}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for duplicate stage name")
	}
}

func TestPipelineRaw_Build_DisabledStageExcluded(t *testing.T) {
	r := &PipelineRaw{
		Stages: []PipelineStageRaw{
			{Name: "filter", Enabled: boolPtrP(true)},
			{Name: "rank", Enabled: boolPtrP(false)},
			{Name: "truncate"},
		},
	}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.EnabledStages) != 2 {
		t.Fatalf("expected 2 enabled stages, got %d", len(cfg.EnabledStages))
	}
	if cfg.EnabledStages[0] != "filter" || cfg.EnabledStages[1] != "truncate" {
		t.Errorf("unexpected stages: %v", cfg.EnabledStages)
	}
}

func TestPipelineRaw_Build_AllEnabled(t *testing.T) {
	r := &PipelineRaw{
		Stages: []PipelineStageRaw{
			{Name: "normalize"},
			{Name: "redact"},
			{Name: "classify"},
		},
	}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.EnabledStages) != 3 {
		t.Errorf("expected 3 stages, got %d", len(cfg.EnabledStages))
	}
}
