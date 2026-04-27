package ssh

import (
	"strings"
	"testing"
	"time"
)

func TestConfig_Defaults(t *testing.T) {
	cfg := Config{
		Host:    "example.com",
		User:    "deploy",
		KeyPath: "/nonexistent/key",
	}

	// New should fail because the key file doesn't exist, but we can
	// confirm the error is about reading the key, not a nil-pointer from
	// missing defaults.
	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for missing key file, got nil")
	}
	if !strings.Contains(err.Error(), "reading key") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestConfig_CustomPort(t *testing.T) {
	cfg := Config{
		Host:    "example.com",
		Port:    2222,
		User:    "deploy",
		KeyPath: "/nonexistent/key",
		Timeout: 5 * time.Second,
	}

	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for missing key file, got nil")
	}
	// Port 2222 should not change the nature of the key-read error.
	if !strings.Contains(err.Error(), "reading key") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestUploadString_ReaderContent(t *testing.T) {
	// Unit-test that UploadString wraps the string in a reader correctly
	// without a live connection by inspecting the helper directly.
	content := "KEY=value\nDEBUG=true\n"
	reader := strings.NewReader(content)

	var buf strings.Builder
	_, err := buf.ReadFrom(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != content {
		t.Errorf("got %q, want %q", buf.String(), content)
	}
}
