package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func annotateResults() []DriftResult {
	return []DriftResult{
		{Key: "db.host", Declared: ptrStr("localhost"), Live: ptrStr("prod-db")},
		{Key: "cache.ttl", Declared: ptrStr("300"), Live: ptrStr("600")},
		{Key: "db.port", Declared: ptrStr("5432"), Live: ptrStr("5433")},
	}
}

func TestAnnotator_NoRules(t *testing.T) {
	a := NewAnnotator(nil)
	results := annotateResults()
	out := a.Annotate(results)
	assert.Len(t, out, 3)
	for _, r := range out {
		assert.Empty(t, r.Tags)
	}
}

func TestAnnotator_PrefixMatch(t *testing.T) {
	a := NewAnnotator([]AnnotationRule{
		{Prefix: "db.", Annotations: []Annotation{{Key: "team", Value: "data"}}},
	})
	out := a.Annotate(annotateResults())
	assert.Contains(t, out[0].Tags, "team=data") // db.host
	assert.Empty(t, out[1].Tags)                 // cache.ttl
	assert.Contains(t, out[2].Tags, "team=data") // db.port
}

func TestAnnotator_MultipleAnnotations(t *testing.T) {
	a := NewAnnotator([]AnnotationRule{
		{
			Prefix: "db.",
			Annotations: []Annotation{
				{Key: "team", Value: "data"},
				{Key: "env", Value: "prod"},
			},
		},
	})
	out := a.Annotate(annotateResults())
	assert.Contains(t, out[0].Tags, "team=data")
	assert.Contains(t, out[0].Tags, "env=prod")
}

func TestAnnotator_DeduplicatesTags(t *testing.T) {
	a := NewAnnotator([]AnnotationRule{
		{Prefix: "db.", Annotations: []Annotation{{Key: "team", Value: "data"}}},
		{Prefix: "db.host", Annotations: []Annotation{{Key: "team", Value: "data"}}},
	})
	out := a.Annotate(annotateResults())
	count := 0
	for _, t := range out[0].Tags {
		if t == "team=data" {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestAnnotator_NoMutation(t *testing.T) {
	a := NewAnnotator([]AnnotationRule{
		{Prefix: "db.", Annotations: []Annotation{{Key: "team", Value: "data"}}},
	})
	original := annotateResults()
	a.Annotate(original)
	assert.Empty(t, original[0].Tags)
}
