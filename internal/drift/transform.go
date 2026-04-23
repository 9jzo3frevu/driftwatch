package drift

import (
	"fmt"
	"strings"
)

// TransformConfig holds configuration for value transformation rules.
type TransformConfig struct {
	Rules []TransformRule
}

// TransformRule defines a key pattern and the transformation to apply.
type TransformRule struct {
	Prefix string
	Op     TransformOp
	Arg    string
}

// TransformOp is the type of transformation to apply.
type TransformOp string

const (
	TransformOpPrefix  TransformOp = "prepend"
	TransformOpSuffix  TransformOp = "append"
	TransformOpReplace TransformOp = "replace"
	TransformOpUpper   TransformOp = "upper"
	TransformOpLower   TransformOp = "lower"
)

// Transformer applies value transformations to drift results.
type Transformer struct {
	cfg TransformConfig
}

// NewTransformer returns a Transformer with the given config.
func NewTransformer(cfg TransformConfig) *Transformer {
	return &Transformer{cfg: cfg}
}

// Apply returns a new slice of DriftResults with declared/live values
// transformed according to matching rules. Original results are not mutated.
func (t *Transformer) Apply(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	for i, r := range results {
		copy := r
		for _, rule := range t.cfg.Rules {
			if !strings.HasPrefix(r.Key, rule.Prefix) {
				continue
			}
			if copy.Declared != nil {
				v := applyOp(*copy.Declared, rule)
				copy.Declared = &v
			}
			if copy.Live != nil {
				v := applyOp(*copy.Live, rule)
				copy.Live = &v
			}
		}
		out[i] = copy
	}
	return out
}

func applyOp(val string, rule TransformRule) string {
	switch rule.Op {
	case TransformOpPrefix:
		return fmt.Sprintf("%s%s", rule.Arg, val)
	case TransformOpSuffix:
		return fmt.Sprintf("%s%s", val, rule.Arg)
	case TransformOpReplace:
		return rule.Arg
	case TransformOpUpper:
		return strings.ToUpper(val)
	case TransformOpLower:
		return strings.ToLower(val)
	default:
		return val
	}
}
