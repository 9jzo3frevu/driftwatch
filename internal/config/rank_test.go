package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRankRaw_Build_Disabled(t *testing.T) {
	cfg, err := (RankRaw{Enabled: false}).Build()
	require.NoError(t, err)
	assert.False(t, cfg.Enabled)
}

func TestRankRaw_Build_DefaultBy(t *testing.T) {
	cfg, err := (RankRaw{Enabled: true}).Build()
	require.NoError(t, err)
	assert.Equal(t, "severity", cfg.By)
}

func TestRankRaw_Build_ValidValues(t *testing.T) {
	for _, by := range []string{"key", "severity", "service"} {
		cfg, err := (RankRaw{Enabled: true, By: by}).Build()
		require.NoError(t, err)
		assert.Equal(t, by, cfg.By)
	}
}

func TestRankRaw_Build_InvalidBy(t *testing.T) {
	_, err := (RankRaw{Enabled: true, By: "timestamp"}).Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rank.by")
}

func TestRankRaw_Build_Descending(t *testing.T) {
	cfg, err := (RankRaw{Enabled: true, Descending: true}).Build()
	require.NoError(t, err)
	assert.True(t, cfg.Descending)
}
