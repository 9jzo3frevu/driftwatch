package drift

import (
	"fmt"
	"regexp"
)

// ValidateConfig holds rules for validating drift result field values.
type ValidateConfig struct {
	Rules []ValidateRule
}

// ValidateRule defines a pattern-based validation rule for a key prefix.
type ValidateRule struct {
	Prefix  string
	Pattern *regexp.Regexp
	Message string
}

// ValidationError is attached to a DriftResult when its live value fails validation.
type ValidationError struct {
	Key     string
	Message string
}

// Validator checks live values in drift results against declared patterns.
type Validator struct {
	cfg ValidateConfig
}

// NewValidator creates a Validator with the provided config.
func NewValidator(cfg ValidateConfig) *Validator {
	return &Validator{cfg: cfg}
}

// Validate applies validation rules to each result and appends a validation
// annotation when a live value does not match the expected pattern.
func (v *Validator) Validate(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	copy(out, results)

	for i, r := range out {
		if r.Live == nil {
			continue
		}
		for _, rule := range v.cfg.Rules {
			if len(r.Key) < len(rule.Prefix) || r.Key[:len(rule.Prefix)] != rule.Prefix {
				continue
			}
			if !rule.Pattern.MatchString(*r.Live) {
				msg := rule.Message
				if msg == "" {
					msg = fmt.Sprintf("value %q does not match pattern %s", *r.Live, rule.Pattern)
				}
				out[i].Annotations = appendUnique(r.Annotations, "validation_error:"+msg)
			}
			break
		}
	}
	return out
}
