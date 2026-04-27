package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Host represents a remote server configuration.
type Host struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Identity string `yaml:"identity"`
}

// Config holds the top-level haul configuration.
type Config struct {
	Version int      `yaml:"version"`
	Hosts   []Host   `yaml:"hosts"`
	Files   []string `yaml:"files"`
}

// DefaultPort is the default SSH port used when none is specified.
const DefaultPort = 22

// Load reads and parses a haul config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	cfg.applyDefaults()
	return &cfg, nil
}

// validate checks that the config contains the minimum required fields.
func (c *Config) validate() error {
	if len(c.Hosts) == 0 {
		return fmt.Errorf("at least one host must be defined")
	}
	for i, h := range c.Hosts {
		if h.Address == "" {
			return fmt.Errorf("host[%d]: address is required", i)
		}
		if h.User == "" {
			return fmt.Errorf("host[%d]: user is required", i)
		}
	}
	if len(c.Files) == 0 {
		return fmt.Errorf("at least one file path must be specified")
	}
	return nil
}

// applyDefaults fills in optional fields with sensible defaults.
func (c *Config) applyDefaults() {
	for i := range c.Hosts {
		if c.Hosts[i].Port == 0 {
			c.Hosts[i].Port = DefaultPort
		}
		if c.Hosts[i].Name == "" {
			c.Hosts[i].Name = c.Hosts[i].Address
		}
	}
}
