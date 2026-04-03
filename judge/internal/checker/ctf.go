package checker

// CTFChecker checks CTF-style submissions by comparing a submitted flag
// against the expected flag. The comparison is case-sensitive and exact.
type CTFChecker struct{}

// Check compares the submitted flag (actual) with the expected flag.
// Returns AC if they match exactly (case-sensitive), WA otherwise.
func (c *CTFChecker) Check(actual string, expected string) Verdict {
	// Trim whitespace from both sides for robustness.
	actual = trimRight(actual)
	expected = trimRight(expected)

	if actual == expected {
		return VerdictAC
	}
	return VerdictWA
}
