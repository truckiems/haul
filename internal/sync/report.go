package sync

import (
	"fmt"
	"io"
)

// PrintReport writes a human-readable summary of sync results to w.
func PrintReport(w io.Writer, results []Result) {
	successCount := 0
	failCount := 0

	for _, r := range results {
		if r.Success {
			successCount++
			fmt.Fprintf(w, "  ✓ %s\n", r.Host)
		} else {
			failCount++
			fmt.Fprintf(w, "  ✗ %s — %v\n", r.Host, r.Err)
		}
	}

	fmt.Fprintf(w, "\nSync complete: %d succeeded, %d failed.\n", successCount, failCount)
}

// HasFailures returns true if any result was not successful.
func HasFailures(results []Result) bool {
	for _, r := range results {
		if !r.Success {
			return true
		}
	}
	return false
}
