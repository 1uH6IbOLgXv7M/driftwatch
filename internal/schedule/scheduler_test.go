package schedule_test

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/schedule"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestScheduler_RunsJobImmediately(t *testing.T) {
	var count atomic.Int32
	job := func(_ context.Context) error {
		count.Add(1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	s := schedule.New(10*time.Second, job, newTestLogger())
	_ = s.Run(ctx)

	if count.Load() < 1 {
		t.Fatal("expected job to run at least once immediately")
	}
}

func TestScheduler_RunsJobOnTick(t *testing.T) {
	var count atomic.Int32
	job := func(_ context.Context) error {
		count.Add(1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	s := schedule.New(100*time.Millisecond, job, newTestLogger())
	_ = s.Run(ctx)

	// immediate run + ~3 ticks within 350 ms
	if count.Load() < 3 {
		t.Fatalf("expected at least 3 runs, got %d", count.Load())
	}
}

func TestScheduler_ContinuesOnJobError(t *testing.T) {
	var count atomic.Int32
	job := func(_ context.Context) error {
		count.Add(1)
		return errors.New("transient error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	s := schedule.New(80*time.Millisecond, job, newTestLogger())
	_ = s.Run(ctx)

	if count.Load() < 2 {
		t.Fatalf("expected scheduler to continue after error, got %d runs", count.Load())
	}
}

func TestScheduler_ReturnsCtxErrOnCancel(t *testing.T) {
	job := func(_ context.Context) error { return nil }

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	s := schedule.New(1*time.Second, job, newTestLogger())
	err := s.Run(ctx)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
