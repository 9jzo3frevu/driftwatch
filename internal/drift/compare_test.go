package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare_NoDrift(t *testing.T) {
	declared := map[string]string{"host": "localhost", "port": "8080"}
	live := map[string]string{"host": "localhost", "port": "8080"}
	results := Compare(declared, live, CompareOptions{Mode: CompareModeExact})
	for _, r := range results {
		assert.False(t, r.Drifted, "expected no drift for key %s", r.Key)
	}
}

func TestCompare_ModifiedValue(t *testing.T) {
	declared := map[string]string{"port": "8080"}
	live := map[string]string{"port": "9090"}
	results := Compare(declared, live, CompareOptions{})
	assert.Len(t, results, 1)
	assert.True(t, results[0].Drifted)
	assert.Contains(t, results[0].Reason, "mismatch")
}

func TestCompare_MissingLiveKey(t *testing.T) {
	declared := map[string]string{"timeout": "30s"}
	live := map[string]string{}
	results := Compare(declared, live, CompareOptions{})
	assert.Len(t, results, 1)
	assert.True(t, results[0].Drifted)
	assert.Equal(t, "key missing from live", results[0].Reason)
}

func TestCompare_ExtraLiveKey_ExactMode(t *testing.T) {
	declared := map[string]string{"host": "localhost"}
	live := map[string]string{"host": "localhost", "extra": "value"}
	results := Compare(declared, live, CompareOptions{Mode: CompareModeExact})
	var drifted []CompareResult
	for _, r := range results {
		if r.Drifted {
			drifted = append(drifted, r)
		}
	}
	assert.Len(t, drifted, 1)
	assert.Equal(t, "extra", drifted[0].Key)
}

func TestCompare_ExtraLiveKey_SubsetMode(t *testing.T) {
	declared := map[string]string{"host": "localhost"}
	live := map[string]string{"host": "localhost", "extra": "value"}
	results := Compare(declared, live, CompareOptions{Mode: CompareModeSubset})
	for _, r := range results {
		assert.False(t, r.Drifted, "subset mode should not flag extra live key")
	}
}

func TestCompare_ServiceField(t *testing.T) {
	opts := CompareOptions{Mode: CompareModeExact, Service: "api"}
	results := Compare(map[string]string{"k": "v"}, map[string]string{"k": "v"}, opts)
	assert.Len(t, results, 1)
	assert.False(t, results[0].Drifted)
}
