package env

// MergeStrategy controls how duplicate keys are handled during a merge.
type MergeStrategy int

const (
	// StrategyOverwrite replaces existing keys with values from the source.
	StrategyOverwrite MergeStrategy = iota
	// StrategyKeepExisting preserves existing keys and only adds new ones.
	StrategyKeepExisting
)

// Merge combines entries from src into dst according to the given strategy.
// It returns a new File; neither dst nor src is mutated.
func Merge(dst, src *File, strategy MergeStrategy) *File {
	existing := make(map[string]int, len(dst.Entries))
	result := make([]Entry, len(dst.Entries))
	copy(result, dst.Entries)

	for i, e := range result {
		existing[e.Key] = i
	}

	for _, e := range src.Entries {
		if idx, found := existing[e.Key]; found {
			if strategy == StrategyOverwrite {
				result[idx] = e
			}
			// StrategyKeepExisting: do nothing
		} else {
			existing[e.Key] = len(result)
			result = append(result, e)
		}
	}

	return &File{Path: dst.Path, Entries: result}
}
