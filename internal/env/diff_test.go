package env

import (
	"testing"
)

func TestDiff_NoDifferences(t *testing.T) {
	local := map[string]string{"FOO": "bar", "BAZ": "qux"}
	remote := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Diff(local, remote)

	if result.HasDiff() {
		t.Errorf("expected no diff, got: %s", result.Summary())
	}
}

func TestDiff_Added(t *testing.T) {
	local := map[string]string{"FOO": "bar", "NEW_KEY": "value"}
	remote := map[string]string{"FOO": "bar"}

	result := Diff(local, remote)

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added key, got %d", len(result.Added))
	}
	if result.Added["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value in added, got %q", result.Added["NEW_KEY"])
	}
	if len(result.Removed) != 0 || len(result.Changed) != 0 {
		t.Errorf("expected no removed or changed, got: %s", result.Summary())
	}
}

func TestDiff_Removed(t *testing.T) {
	local := map[string]string{"FOO": "bar"}
	remote := map[string]string{"FOO": "bar", "OLD_KEY": "old"}

	result := Diff(local, remote)

	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed key, got %d", len(result.Removed))
	}
	if result.Removed["OLD_KEY"] != "old" {
		t.Errorf("expected OLD_KEY=old in removed")
	}
}

func TestDiff_Changed(t *testing.T) {
	local := map[string]string{"FOO": "new_value"}
	remote := map[string]string{"FOO": "old_value"}

	result := Diff(local, remote)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(result.Changed))
	}
	pair, ok := result.Changed["FOO"]
	if !ok {
		t.Fatal("expected FOO in changed")
	}
	if pair[0] != "new_value" || pair[1] != "old_value" {
		t.Errorf("expected [new_value old_value], got %v", pair)
	}
}

func TestDiff_Summary(t *testing.T) {
	local := map[string]string{"A": "1", "B": "changed"}
	remote := map[string]string{"B": "original", "C": "removed"}

	result := Diff(local, remote)
	summary := result.Summary()

	if summary == "" {
		t.Error("expected non-empty summary")
	}
	// +1 added (A), -1 removed (C), ~1 changed (B)
	expected := "+1 added, -1 removed, ~1 changed"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
