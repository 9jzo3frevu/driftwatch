package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagRuleRaw_Build_MissingPrefix(t *testing.T) {
	raw := TagRuleRaw{Tags: map[string]string{"team": "data"}}
	_, err := raw.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing prefix")
}

func TestTagRuleRaw_Build_NoTags(t *testing.T) {
	raw := TagRuleRaw{Prefix: "db."}
	_, err := raw.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no tags")
}

func TestTagRuleRaw_Build_Valid(t *testing.T) {
	raw := TagRuleRaw{
		Prefix: "db.",
		Tags:   map[string]string{"team": "data"},
	}
	rule, err := raw.Build()
	require.NoError(t, err)
	assert.Equal(t, "db.", rule.Prefix)
	assert.Len(t, rule.Tags, 1)
	assert.Equal(t, "team", rule.Tags[0].Key)
	assert.Equal(t, "data", rule.Tags[0].Value)
}

func TestBuildTagRules_Valid(t *testing.T) {
	raws := []TagRuleRaw{
		{Prefix: "db.", Tags: map[string]string{"team": "data"}},
		{Prefix: "cache.", Tags: map[string]string{"env": "prod"}},
	}
	rules, err := BuildTagRules(raws)
	require.NoError(t, err)
	assert.Len(t, rules, 2)
}

func TestBuildTagRules_PropagatesError(t *testing.T) {
	raws := []TagRuleRaw{
		{Prefix: "db.", Tags: map[string]string{"team": "data"}},
		{Prefix: "", Tags: map[string]string{"env": "prod"}},
	}
	_, err := BuildTagRules(raws)
	require.Error(t, err)
}
