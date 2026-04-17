package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func tagResults() []DriftResult {
	return []DriftResult{
		{Key: "db.host", Declared: ptrStr("localhost"), Live: ptrStr("prod-db")},
		{Key: "cache.ttl", Declared: ptrStr("300"), Live: ptrStr("600")},
		{Key: "app.name", Declared: ptrStr("svc"), Live: ptrStr("svc")},
	}
}

func TestTagger_NoRules(t *testing.T) {
	tagger := NewTagger(nil)
	results := tagger.Apply(tagResults())
	for _, r := range results {
		assert.Empty(t, r.Tags)
	}
}

func TestTagger_PrefixMatch(t *testing.T) {
	rules := []TagRule{
		{Prefix: "db.", Tags: []Tag{{Key: "team", Value: "data"}}},
	}
	tagger := NewTagger(rules)
	results := tagger.Apply(tagResults())

	assert.Equal(t, []Tag{{Key: "team", Value: "data"}}, results[0].Tags)
	assert.Empty(t, results[1].Tags)
	assert.Empty(t, results[2].Tags)
}

func TestTagger_MultipleRulesMatch(t *testing.T) {
	rules := []TagRule{
		{Prefix: "db.", Tags: []Tag{{Key: "team", Value: "data"}}},
		{Prefix: "db.", Tags: []Tag{{Key: "env", Value: "prod"}}},
	}
	tagger := NewTagger(rules)
	results := tagger.Apply(tagResults())

	assert.Len(t, results[0].Tags, 2)
}

func TestTagger_NoMutation(t *testing.T) {
	original := tagResults()
	rules := []TagRule{
		{Prefix: "db.", Tags: []Tag{{Key: "team", Value: "data"}}},
	}
	tagger := NewTagger(rules)
	tagger.Apply(original)
	assert.Empty(t, original[0].Tags)
}
