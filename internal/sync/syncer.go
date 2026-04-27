package sync

import (
	"fmt"
	"os"

	"github.com/haul/internal/config"
	"github.com/haul/internal/ssh"
)

// Result holds the outcome of a sync operation for a single host.
type Result struct {
	Host    string
	Success bool
	Err     error
}

// Syncer orchestrates uploading env files to remote hosts.
type Syncer struct {
	cfg *config.Config
}

// New creates a new Syncer from the given config.
func New(cfg *config.Config) *Syncer {
	return &Syncer{cfg: cfg}
}

// Run syncs the configured env file to all hosts concurrently.
func (s *Syncer) Run() []Result {
	results := make(chan Result, len(s.cfg.Hosts))

	for _, host := range s.cfg.Hosts {
		go func(h config.Host) {
			err := s.syncHost(h)
			results <- Result{Host: h.Address, Success: err == nil, Err: err}
		}(host)
	}

	out := make([]Result, 0, len(s.cfg.Hosts))
	for range s.cfg.Hosts {
		out = append(out, <-results)
	}
	return out
}

func (s *Syncer) syncHost(host config.Host) error {
	client := ssh.New(ssh.Config{
		Address:    host.Address,
		User:       host.User,
		KeyPath:    host.KeyPath,
		Port:       host.Port,
	})

	f, err := os.Open(s.cfg.EnvFile)
	if err != nil {
		return fmt.Errorf("open env file: %w", err)
	}
	defer f.Close()

	if err := client.UploadReader(f, host.RemotePath); err != nil {
		return fmt.Errorf("upload to %s: %w", host.Address, err)
	}
	return nil
}
