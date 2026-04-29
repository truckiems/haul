package env

// MergeStrategy defines how conflicting keys are handled during a merge.
type MergeStrategy int

const (
	// StrategyOverwrite replaces existing keys with incoming values.
	StrategyOverwrite MergeStrategy = iota
	// StrategyKeepExisting preserves existing keys and only adds new ones.
	StrategyKeepExisting
)

// Merge combines base and override maps according to the given strategy.
// The base map is never mutated; a new map is returned.
func Merge(base, override map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	for k, v := range override {
		switch strategy {
		case StrategyOverwrite:
			result[k] = v
		case StrategyKeepExisting:
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}

	return result
}

// MergeAll merges a slice of env maps left-to-right using the given strategy.
// Earlier maps take precedence with StrategyKeepExisting; later maps win with StrategyOverwrite.
func MergeAll(maps []map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		result = Merge(result, m, strategy)
	}
	return result
}
