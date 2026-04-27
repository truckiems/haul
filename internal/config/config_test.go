package config

import (
	"os"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "haul-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
version: 1
hosts:
  - address: 192.168.1.10
    user: deploy
files:
  - /etc/app/.env
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(cfg.Hosts))
	}
	if cfg.Hosts[0].Port != DefaultPort {
		t.Errorf("expected default port %d, got %d", DefaultPort, cfg.Hosts[0].Port)
	}
	if cfg.Hosts[0].Name != "192.168.1.10" {
		t.Errorf("expected name to default to address, got %q", cfg.Hosts[0].Name)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	path := writeTempConfig(t, `
hosts:
  - user: deploy
files:
  - /etc/app/.env
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing address, got nil")
	}
}

func TestLoad_NoFiles(t *testing.T) {
	path := writeTempConfig(t, `
hosts:
  - address: 10.0.0.1
    user: admin
files: []
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty files list, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/haul.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
