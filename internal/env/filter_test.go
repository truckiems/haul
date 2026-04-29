package env

import (
	"testing"
)

func TestFilter_NoPrefixesNoExcludes(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "DB_URL": "postgres://", "SECRET": "abc"}
	got := Filter(env, FilterOptions{})
	if len(got) != 3 {
		t.Errorf("expected 3 keys, got %d", len(got))
	}
}

func TestFilter_WithPrefix(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	got := Filter(env, FilterOptions{Prefixes: []string{"APP_"}})
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
	if _, ok := got["DB_URL"]; ok {
		t.Error("DB_URL should have been filtered out")
	}
}

func TestFilter_WithExclude(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "SECRET": "abc", "DB_URL": "postgres://"}
	got := Filter(env, FilterOptions{Exclude: []string{"SECRET"}})
	if _, ok := got["SECRET"]; ok {
		t.Error("SECRET should have been excluded")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestFilter_PrefixAndExclude(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_SECRET": "xyz", "DB_URL": "postgres://"}
	got := Filter(env, FilterOptions{
		Prefixes: []string{"APP_"},
		Exclude:  []string{"APP_SECRET"},
	})
	if len(got) != 1 {
		t.Errorf("expected 1 key, got %d", len(got))
	}
	if _, ok := got["APP_HOST"]; !ok {
		t.Error("APP_HOST should be present")
	}
}

func TestFilter_MultiplePrefixes(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "DB_URL": "postgres://", "LOG_LEVEL": "info", "SECRET": "abc"}
	got := Filter(env, FilterOptions{Prefixes: []string{"APP_", "DB_"}})
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestKeys_Sorted(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := Keys(env)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], k)
		}
	}
}
