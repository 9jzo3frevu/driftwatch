package drift

import "fmt"

// Explanation holds a human-readable explanation for a single drift result.
type Explanation struct {
	Key     string
	Kind    string
	Message string
	Hint    string
}

// Explainer generates explanations for drift results.
type Explainer struct {
	hints map[string]string
}

// NewExplainer returns an Explainer with optional per-key hints.
func NewExplainer(hints map[string]string) *Explainer {
	if hints == nil {
		hints = map[string]string{}
	}
	return &Explainer{hints: hints}
}

// Explain returns an Explanation slice for the given results.
func (e *Explainer) Explain(results []DriftResult) []Explanation {
	out := make([]Explanation, 0, len(results))
	for _, r := range results {
		out = append(out, e.explainOne(r))
	}
	return out
}

func (e *Explainer) explainOne(r DriftResult) Explanation {
	var msg string
	switch r.Kind {
	case KindModified:
		msg = fmt.Sprintf("key %q changed from %q to %q", r.Key, ptrStr(r.Declared), ptrStr(r.Live))
	case KindMissing:
		msg = fmt.Sprintf("key %q is declared but absent in live config", r.Key)
	case KindExtra:
		msg = fmt.Sprintf("key %q exists in live config but is not declared", r.Key)
	default:
		msg = fmt.Sprintf("key %q has unknown drift kind %q", r.Key, r.Kind)
	}
	hint := e.hints[r.Key]
	return Explanation{
		Key:     r.Key,
		Kind:    string(r.Kind),
		Message: msg,
		Hint:    hint,
	}
}
