package env

import "fmt"

// DiffResult holds the differences between two env maps.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [local, remote]
}

// Diff compares a local env map against a remote env map and returns the
// differences. local is the source of truth.
func Diff(local, remote map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	for k, localVal := range local {
		remoteVal, exists := remote[k]
		if !exists {
			result.Added[k] = localVal
		} else if localVal != remoteVal {
			result.Changed[k] = [2]string{localVal, remoteVal}
		}
	}

	for k, remoteVal := range remote {
		if _, exists := local[k]; !exists {
			result.Removed[k] = remoteVal
		}
	}

	return result
}

// HasDiff returns true if there are any differences in the result.
func (d DiffResult) HasDiff() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Summary returns a human-readable summary of the diff.
func (d DiffResult) Summary() string {
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed",
		len(d.Added), len(d.Removed), len(d.Changed))
}
