package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/yourorg/driftwatch/internal/config"
)

const defaultConfigPath = "driftwatch.yaml"

func main() {
	cfgPath := flag.String("config", defaultConfigPath, "path to driftwatch config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "driftwatch: failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := newLogger(cfg.Daemon.LogLevel)
	logger.Info("driftwatch starting",
		"provider", cfg.Cloud.Provider,
		"region", cfg.Cloud.Region,
		"poll_interval", cfg.Daemon.PollInterval,
		"workspace", cfg.Terraform.Workspace,
	)

	// TODO: initialise drift detector and begin polling loop.
	logger.Info("daemon initialised — polling not yet implemented")
}

func newLogger(level string) *slog.Logger {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}
