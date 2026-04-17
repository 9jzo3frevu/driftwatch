package config

// ScoreRaw holds configuration for drift scoring thresholds.
type ScoreRaw struct {
	WarnThreshold *int `yaml:"warn_threshold"`
	FailThreshold *int `yaml:"fail_threshold"`
}

// ScoreConfig holds validated scoring thresholds.
type ScoreConfig struct {
	WarnThreshold int
	FailThreshold int
}

const defaultWarnThreshold = 10
const defaultFailThreshold = 100

// Build validates and returns a ScoreConfig.
func (r ScoreRaw) Build() ScoreConfig {
	warn := defaultWarnThreshold
	if r.WarnThreshold != nil {
		warn = *r.WarnThreshold
	}
	fail := defaultFailThreshold
	if r.FailThreshold != nil {
		fail = *r.FailThreshold
	}
	return ScoreConfig{
		WarnThreshold: warn,
		FailThreshold: fail,
	}
}
