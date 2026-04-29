package env

import (
	"testing"
)

func TestMerge_Overwrite(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"B": "99", "C": "3"}

	result := Merge(base, override, StrategyOverwrite)

	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
	if result["B"] != "99" {
		t.Errorf("expected B=99, got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestMerge_KeepExisting(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"B": "99", "C": "3"}

	result := Merge(base, override, StrategyKeepExisting)

	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
	if result["B"] != "2" {
		t.Errorf("expected B=2 (preserved), got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestMerge_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"A": "original"}
	override := map[string]string{"A": "changed"}

	_ = Merge(base, override, StrategyOverwrite)

	if base["A"] != "original" {
		t.Errorf("base map was mutated: got %s", base["A"])
	}
}

func TestMergeAll_Overwrite(t *testing.T) {
	maps := []map[string]string{
		{"A": "1", "B": "2"},
		{"B": "3", "C": "4"},
		{"C": "99"},
	}

	result := MergeAll(maps, StrategyOverwrite)

	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
	if result["B"] != "3" {
		t.Errorf("expected B=3, got %s", result["B"])
	}
	if result["C"] != "99" {
		t.Errorf("expected C=99, got %s", result["C"])
	}
}

func TestMergeAll_KeepExisting(t *testing.T) {
	maps := []map[string]string{
		{"A": "first"},
		{"A": "second", "B": "only"},
	}

	result := MergeAll(maps, StrategyKeepExisting)

	if result["A"] != "first" {
		t.Errorf("expected A=first, got %s", result["A"])
	}
	if result["B"] != "only" {
		t.Errorf("expected B=only, got %s", result["B"])
	}
}

func TestMergeAll_Empty(t *testing.T) {
	result := MergeAll([]map[string]string{}, StrategyOverwrite)
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}
