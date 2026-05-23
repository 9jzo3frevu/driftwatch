package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func partitionResults() []DriftResult {
	return []DriftResult{
		{Key: "db.host", Service: "api", Drifted: true, Severity: SeverityHigh},
		{Key: "db.port", Service: "api", Drifted: true, Severity: SeverityLow},
		{Key: "cache.host", Service: "worker", Drifted: true, Severity: SeverityHigh},
		{Key: "cache.ttl", Service: "worker", Drifted: false, Severity: SeverityLow},
	}
}

func TestPartitionResults_ByService(t *testing.T) {
	cfg := DefaultPartitionConfig()
	parts := PartitionResults(partitionResults(), cfg)

	assert.Len(t, parts, 2)
	labels := map[string]int{}
	for _, p := range parts {
		labels[p.Label] = len(p.Results)
	}
	assert.Equal(t, 2, labels["api"])
	assert.Equal(t, 2, labels["worker"])
}

func TestPartitionResults_BySeverity(t *testing.T) {
	cfg := PartitionConfig{By: "severity"}
	parts := PartitionResults(partitionResults(), cfg)

	assert.Len(t, parts, 2)
	labels := map[string]int{}
	for _, p := range parts {
		labels[p.Label] = len(p.Results)
	}
	assert.Equal(t, 2, labels[string(SeverityHigh)])
	assert.Equal(t, 2, labels[string(SeverityLow)])
}

func TestPartitionResults_ByKeyPrefix(t *testing.T) {
	cfg := PartitionConfig{By: "key_prefix"}
	parts := PartitionResults(partitionResults(), cfg)

	assert.Len(t, parts, 2)
	labels := map[string]bool{}
	for _, p := range parts {
		labels[p.Label] = true
	}
	assert.True(t, labels["db"])
	assert.True(t, labels["cache"])
}

func TestPartitionResults_MaxSize(t *testing.T) {
	cfg := PartitionConfig{By: "service", MaxSize: 1}
	parts := PartitionResults(partitionResults(), cfg)

	for _, p := range parts {
		assert.LessOrEqual(t, len(p.Results), 1)
	}
}

func TestPartitionResults_Empty(t *testing.T) {
	parts := PartitionResults(nil, DefaultPartitionConfig())
	assert.Nil(t, parts)
}

func TestPartitionResults_UnknownService(t *testing.T) {
	results := []DriftResult{
		{Key: "x", Service: "", Drifted: true, Severity: SeverityLow},
	}
	parts := PartitionResults(results, DefaultPartitionConfig())
	assert.Len(t, parts, 1)
	assert.Equal(t, "unknown", parts[0].Label)
}
