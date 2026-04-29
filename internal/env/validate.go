package env

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation issue found in an env map.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) Summary() string {
	if !r.HasErrors() {
		return "all keys valid"
	}
	lines := make([]string, 0, len(r.Errors))
	for _, e := range r.Errors {
		lines = append(lines, "  - "+e.Error())
	}
	return fmt.Sprintf("%d validation error(s):\n%s", len(r.Errors), strings.Join(lines, "\n"))
}

// Validate checks that all keys in the env map are non-empty, contain only
// allowed characters (A-Z, a-z, 0-9, _), and do not start with a digit.
func Validate(env map[string]string) ValidationResult {
	result := ValidationResult{}
	for k := range env {
		if k == "" {
			result.Errors = append(result.Errors, ValidationError{Key: k, Message: "key must not be empty"})
			continue
		}
		if k[0] >= '0' && k[0] <= '9' {
			result.Errors = append(result.Errors, ValidationError{Key: k, Message: "key must not start with a digit"})
		}
		for _, ch := range k {
			if !isValidKeyChar(ch) {
				result.Errors = append(result.Errors, ValidationError{
					Key:     k,
					Message: fmt.Sprintf("invalid character %q in key", ch),
				})
				break
			}
		}
	}
	return result
}

func isValidKeyChar(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '_'
}
