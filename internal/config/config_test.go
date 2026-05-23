package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*.yaml")
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
	raw := `
daemon:
  poll_interval: 10m
  log_level: debug
cloud:
  provider: aws
  region: us-east-1
  profile: default
terraform:
  state_file: terraform.tfstate
  workspace: staging
`
	path := writeTempConfig(t, raw)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Daemon.PollInterval != 10*time.Minute {
		t.Errorf("poll_interval: got %v, want 10m", cfg.Daemon.PollInterval)
	}
	if cfg.Cloud.Provider != "aws" {
		t.Errorf("provider: got %q, want aws", cfg.Cloud.Provider)
	}
	if cfg.Terraform.Workspace != "staging" {
		t.Errorf("workspace: got %q, want staging", cfg.Terraform.Workspace)
	}
}

func TestLoad_Defaults(t *testing.T) {
	raw := `
cloud:
  provider: gcp
  region: us-central1
terraform:
  state_file: terraform.tfstate
`
	path := writeTempConfig(t, raw)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Daemon.PollInterval != 5*time.Minute {
		t.Errorf("default poll_interval: got %v, want 5m", cfg.Daemon.PollInterval)
	}
	if cfg.Daemon.LogLevel != "info" {
		t.Errorf("default log_level: got %q, want info", cfg.Daemon.LogLevel)
	}
	if cfg.Terraform.Workspace != "default" {
		t.Errorf("default workspace: got %q, want default", cfg.Terraform.Workspace)
	}
}

func TestLoad_MissingProvider(t *testing.T) {
	raw := `
cloud:
  region: us-east-1
terraform:
  state_file: terraform.tfstate
`
	path := writeTempConfig(t, raw)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestLoad_MissingStateSource(t *testing.T) {
	raw := `
cloud:
  provider: aws
  region: us-east-1
terraform:
  workspace: default
`
	path := writeTempConfig(t, raw)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
