package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func rankResults() []DriftResult {
	return []DriftResult{
		{Key: "db.host", Service: "api", Severity: "high"},
		{Key: "app.port", Service: "worker", Severity: "low"},
		{Key: "cache.ttl", Service: "api", Severity: "critical"},
		{Key: "auth.secret", Service: "auth", Severity: "medium"},
	}
}

func TestRankResults_ByKey_Ascending(t *testing.T) {
	out := RankResults(rankResults(), RankOptions{By: "key"})
	assert.Equal(t, "app.port", out[0].Key)
	assert.Equal(t, "auth.secret", out[1].Key)
	assert.Equal(t, "cache.ttl", out[2].Key)
	assert.Equal(t, "db.host", out[3].Key)
}

func TestRankResults_BySeverity_Ascending(t *testing.T) {
	out := RankResults(rankResults(), RankOptions{By: "severity"})
	assert.Equal(t, "low", out[0].Severity)
	assert.Equal(t, "medium", out[1].Severity)
	assert.Equal(t, "high", out[2].Severity)
	assert.Equal(t, "critical", out[3].Severity)
}

func TestRankResults_BySeverity_Descending(t *testing.T) {
	out := RankResults(rankResults(), RankOptions{By: "severity", Descending: true})
	assert.Equal(t, "critical", out[0].Severity)
	assert.Equal(t, "low", out[3].Severity)
}

func TestRankResults_ByService(t *testing.T) {
	out := RankResults(rankResults(), RankOptions{By: "service"})
	assert.Equal(t, "api", out[0].Service)
	assert.Equal(t, "api", out[1].Service)
	assert.Equal(t, "auth", out[2].Service)
	assert.Equal(t, "worker", out[3].Service)
}

func TestRankResults_NoMutation(t *testing.T) {
	orig := rankResults()
	RankResults(orig, RankOptions{By: "key"})
	assert.Equal(t, "db.host", orig[0].Key)
}
