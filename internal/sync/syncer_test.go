package sync_test

import (
	"os"
	"testing"

	"github.com/haul/internal/config"
	"github.com/haul/internal/sync"
)

func TestNew_ReturnsSyncer(t *testing.T) {
	cfg := &config.Config{
		EnvFile: ".env",
		Hosts:   []config.Host{},
	}
	s := sync.New(cfg)
	if s == nil {
		t.Fatal("expected non-nil syncer")
	}
}

func TestRun_NoHosts(t *testing.T) {
	cfg := &config.Config{
		EnvFile: ".env",
		Hosts:   []config.Host{},
	}
	s := sync.New(cfg)
	results := s.Run()
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRun_MissingEnvFile(t *testing.T) {
	cfg := &config.Config{
		EnvFile: "/nonexistent/.env.missing",
		Hosts: []config.Host{
			{Address: "192.0.2.1", User: "deploy", RemotePath: "/app/.env"},
		},
	}
	s := sync.New(cfg)
	results := s.Run()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Success {
		t.Error("expected failure for missing env file")
	}
	if results[0].Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestRun_ResultsMatchHosts(t *testing.T) {
	f, err := os.CreateTemp("", ".env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("KEY=value\n")
	f.Close()

	cfg := &config.Config{
		EnvFile: f.Name(),
		Hosts: []config.Host{
			{Address: "192.0.2.1", User: "deploy", RemotePath: "/app/.env"},
			{Address: "192.0.2.2", User: "deploy", RemotePath: "/app/.env"},
		},
	}
	s := sync.New(cfg)
	results := s.Run()
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}
