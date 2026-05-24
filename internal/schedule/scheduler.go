// Package schedule provides periodic drift-check scheduling.
package schedule

import (
	"context"
	"log/slog"
	"time"
)

// Job is a function that performs a single drift-check cycle.
type Job func(ctx context.Context) error

// Scheduler runs a Job at a fixed interval until the context is cancelled.
type Scheduler struct {
	interval time.Duration
	job      Job
	logger   *slog.Logger
}

// New creates a Scheduler that will invoke job every interval.
func New(interval time.Duration, job Job, logger *slog.Logger) *Scheduler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Scheduler{
		interval: interval,
		job:      job,
		logger:   logger,
	}
}

// Run starts the scheduling loop. It executes the job immediately on first
// tick and then once every interval. It blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	s.logger.Info("scheduler started", "interval", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run immediately before waiting for the first tick.
	if err := s.runJob(ctx); err != nil {
		s.logger.Error("drift check failed", "error", err)
	}

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := s.runJob(ctx); err != nil {
				s.logger.Error("drift check failed", "error", err)
			}
		}
	}
}

func (s *Scheduler) runJob(ctx context.Context) error {
	start := time.Now()
	s.logger.Info("drift check starting")
	err := s.job(ctx)
	s.logger.Info("drift check finished", "duration", time.Since(start))
	return err
}
