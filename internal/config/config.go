package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level driftwatch daemon configuration.
type Config struct {
	Daemon  DaemonConfig  `yaml:"daemon"`
	Cloud   CloudConfig   `yaml:"cloud"`
	Terraform TerraformConfig `yaml:"terraform"`
}

// DaemonConfig controls polling and logging behaviour.
type DaemonConfig struct {
	PollInterval time.Duration `yaml:"poll_interval"`
	LogLevel     string        `yaml:"log_level"`
}

// CloudConfig holds provider credentials and region.
type CloudConfig struct {
	Provider string `yaml:"provider"`
	Region   string `yaml:"region"`
	Profile  string `yaml:"profile"`
}

// TerraformConfig points driftwatch at a state source.
type TerraformConfig struct {
	StateFile   string `yaml:"state_file"`
	Workspace   string `yaml:"workspace"`
	BackendType string `yaml:"backend_type"`
}

// Load reads a YAML config file from path and returns a validated Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing YAML: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	setDefaults(&cfg)
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Cloud.Provider == "" {
		return fmt.Errorf("config: cloud.provider is required")
	}
	if c.Cloud.Region == "" {
		return fmt.Errorf("config: cloud.region is required")
	}
	if c.Terraform.StateFile == "" && c.Terraform.BackendType == "" {
		return fmt.Errorf("config: one of terraform.state_file or terraform.backend_type is required")
	}
	return nil
}

func setDefaults(c *Config) {
	if c.Daemon.PollInterval == 0 {
		c.Daemon.PollInterval = 5 * time.Minute
	}
	if c.Daemon.LogLevel == "" {
		c.Daemon.LogLevel = "info"
	}
	if c.Terraform.Workspace == "" {
		c.Terraform.Workspace = "default"
	}
}
