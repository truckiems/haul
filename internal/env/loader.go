package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from an env file.
type Entry struct {
	Key   string
	Value string
}

// File holds the parsed contents of an env file.
type File struct {
	Path    string
	Entries []Entry
}

// Load reads and parses a .env file from the given path.
// Lines starting with '#' are treated as comments and skipped.
// Empty lines are also skipped.
func Load(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("env: open %q: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("env: %q line %d: invalid format (expected KEY=VALUE)", path, lineNum)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"`)

		if key == "" {
			return nil, fmt.Errorf("env: %q line %d: empty key", path, lineNum)
		}

		entries = append(entries, Entry{Key: key, Value: value})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env: scan %q: %w", path, err)
	}

	return &File{Path: path, Entries: entries}, nil
}

// String serialises the env file back to its canonical text representation.
func (f *File) String() string {
	var sb strings.Builder
	for _, e := range f.Entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}
