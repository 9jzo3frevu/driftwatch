package drift

// ChangeType describes the kind of configuration drift detected.
type ChangeType string

const (
	// Modified indicates a key exists in both declared and actual state
	// but the values differ.
	Modified ChangeType = "modified"

	// Added indicates a key is present in the actual state but absent
	// from the declared (expected) configuration.
	Added ChangeType = "added"

	// Removed indicates a key is declared in the expected configuration
	// but missing from the actual deployed state.
	Removed ChangeType = "removed"
)

// DriftResult represents a single detected configuration discrepancy.
type DriftResult struct {
	// Key is the configuration key that differs.
	Key string

	// ChangeType classifies the nature of the drift.
	ChangeType ChangeType

	// Expected holds the declared value, or nil when the key was not declared.
	Expected *string

	// Actual holds the deployed value, or nil when the key is absent in
	// the running service.
	Actual *string
}

// Summary returns a human-readable one-line description of the drift result.
func (d DriftResult) Summary() string {
	switch d.ChangeType {
	case Modified:
		return d.Key + ": expected " + ptrStr(d.Expected) + ", got " + ptrStr(d.Actual)
	case Added:
		return d.Key + ": unexpected key with value " + ptrStr(d.Actual)
	case Removed:
		return d.Key + ": declared value " + ptrStr(d.Expected) + " is missing"
	default:
		return d.Key + ": unknown drift type"
	}
}

func ptrStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
