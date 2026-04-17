package drift

import (
	"fmt"
	"strings"
)

// RemediationHint describes a suggested fix for a drift result.
type RemediationHint struct {
	Key        string
	Suggestion string
}

// Remediator generates remediation hints for drift results.
type Remediator struct {
	template string
}

// NewRemediator creates a Remediator with an optional command template.
// The template may contain {key}, {declared}, {live} placeholders.
func NewRemediator(template string) *Remediator {
	if template == "" {
		template = "update {key} from '{live}' to '{declared}'"
	}
	return &Remediator{template: template}
}

// Hints returns a RemediationHint for each actionable DriftResult.
func (r *Remediator) Hints(results []DriftResult) []RemediationHint {
	var hints []RemediationHint
	for _, res := range results {
		if res.Status == StatusMatch {
			continue
		}
		hints = append(hints, RemediationHint{
			Key:        res.Key,
			Suggestion: r.render(res),
		})
	}
	return hints
}

func (r *Remediator) render(res DriftResult) string {
	s := r.template
	s = strings.ReplaceAll(s, "{key}", res.Key)
	s = strings.ReplaceAll(s, "{declared}", fmt.Sprintf("%v", ptrStr(res.Declared)))
	s = strings.ReplaceAll(s, "{live}", fmt.Sprintf("%v", ptrStr(res.Live)))
	return s
}
