package checker

// Verdict represents the result of checking a submission against expected output.
type Verdict string

const (
	// VerdictAC means Accepted - output matches exactly.
	VerdictAC Verdict = "AC"
	// VerdictWA means Wrong Answer - output does not match.
	VerdictWA Verdict = "WA"
	// VerdictPE means Presentation Error - output matches after ignoring whitespace.
	VerdictPE Verdict = "PE"
	// VerdictCE means Compilation Error.
	VerdictCE Verdict = "CE"
	// VerdictRE means Runtime Error.
	VerdictRE Verdict = "RE"
	// VerdictTLE means Time Limit Exceeded.
	VerdictTLE Verdict = "TLE"
	// VerdictMLE means Memory Limit Exceeded.
	VerdictMLE Verdict = "MLE"
)

// Checker defines the interface for output checkers.
type Checker interface {
	// Check compares the actual output with the expected output and returns a verdict.
	Check(actual string, expected string) Verdict
}

// StandardChecker compares stdout with expected output, ignoring trailing whitespace.
type StandardChecker struct{}

// Check performs a standard comparison:
// 1. If outputs are exactly equal -> AC
// 2. If outputs are equal after trimming trailing whitespace on each line -> PE
// 3. Otherwise -> WA
func (s *StandardChecker) Check(actual string, expected string) Verdict {
	if actual == expected {
		return VerdictAC
	}

	// Try trimming trailing whitespace on each line.
	actualTrimmed := trimTrailingWhitespace(actual)
	expectedTrimmed := trimTrailingWhitespace(expected)

	if actualTrimmed == expectedTrimmed {
		return VerdictPE
	}

	return VerdictWA
}

// trimTrailingWhitespace removes trailing whitespace from each line and
// removes trailing newlines from the entire string.
func trimTrailingWhitespace(s string) string {
	lines := splitLines(s)
	for i, line := range lines {
		lines[i] = trimRight(line)
	}
	// Remove trailing empty lines.
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

// splitLines splits a string into lines, handling both \n and \r\n.
func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	// Normalize line endings.
	s = normalizeLineEndings(s)
	lines := splitString(s, '\n')
	return lines
}

// normalizeLineEndings replaces \r\n with \n.
func normalizeLineEndings(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\r' && i+1 < len(s) && s[i+1] == '\n' {
			result = append(result, '\n')
			i++
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}

// trimRight removes trailing whitespace characters from a string.
func trimRight(s string) string {
	i := len(s) - 1
	for i >= 0 && (s[i] == ' ' || s[i] == '\t' || s[i] == '\r' || s[i] == '\n') {
		i--
	}
	return s[:i+1]
}

// splitString splits a string by a separator byte.
func splitString(s string, sep byte) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}
