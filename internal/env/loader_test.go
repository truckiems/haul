package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestLoad_ValidFile(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\nDB_PORT=5432\n")

	file, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(file.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(file.Entries))
	}
	if file.Entries[0].Key != "APP_ENV" || file.Entries[0].Value != "production" {
		t.Errorf("unexpected first entry: %+v", file.Entries[0])
	}
}

func TestLoad_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nFOO=bar\n")

	file, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(file.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(file.Entries))
	}
}

func TestLoad_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="mysecretvalue"` + "\n")

	file, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file.Entries[0].Value != "mysecretvalue" {
		t.Errorf("expected unquoted value, got %q", file.Entries[0].Value)
	}
}

func TestLoad_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestFile_String(t *testing.T) {
	file := &File{
		Entries: []Entry{
			{Key: "FOO", Value: "bar"},
			{Key: "BAZ", Value: "qux"},
		},
	}
	got := file.String()
	want := "FOO=bar\nBAZ=qux\n"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
