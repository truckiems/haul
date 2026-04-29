package env

import "strings"

// FilterOptions controls which keys are included or excluded.
type FilterOptions struct {
	// Only include keys with these prefixes (empty means include all)
	Prefixes []string
	// Exclude keys matching these exact names
	Exclude []string
}

// Filter returns a subset of the env map based on the provided options.
// If opts.Prefixes is non-empty, only keys matching at least one prefix are kept.
// Keys listed in opts.Exclude are always removed.
func Filter(env map[string]string, opts FilterOptions) map[string]string {
	excludeSet := make(map[string]struct{}, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excludeSet[k] = struct{}{}
	}

	result := make(map[string]string)
	for k, v := range env {
		if _, excluded := excludeSet[k]; excluded {
			continue
		}
		if len(opts.Prefixes) == 0 {
			result[k] = v
			continue
		}
		for _, prefix := range opts.Prefixes {
			if strings.HasPrefix(k, prefix) {
				result[k] = v
				break
			}
		}
	}
	return result
}

// Keys returns a sorted slice of keys from the env map.
func Keys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

// sortStrings sorts a slice of strings in place (simple insertion sort for small slices).
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
