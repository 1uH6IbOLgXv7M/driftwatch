package drift_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/your-org/driftwatch/internal/drift"
	"github.com/your-org/driftwatch/internal/state"
)

// stubFetcher is a test double for drift.LiveFetcher.
type stubFetcher struct {
	attrs map[string]map[string]string // keyed by resource name
	err   error
}

func (s *stubFetcher) Fetch(_ context.Context, r state.Resource) (map[string]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.attrs[r.Name], nil
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestDetector_NoDrift(t *testing.T) {
	resources := []state.Resource{
		{Name: "aws_instance.web", Type: "aws_instance", Provider: "aws",
			Attributes: map[string]string{"instance_type": "t3.micro", "ami": "ami-123"}},
	}
	fetcher := &stubFetcher{
		attrs: map[string]map[string]string{
			"aws_instance.web": {"instance_type": "t3.micro", "ami": "ami-123"},
		},
	}

	detector := drift.NewDetector(fetcher, newTestLogger())
	results, err := detector.Check(context.Background(), resources)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Errorf("expected no drift, but drift was reported for keys: %v", results[0].DriftedKeys)
	}
}

func TestDetector_WithDrift(t *testing.T) {
	resources := []state.Resource{
		{Name: "aws_instance.web", Type: "aws_instance", Provider: "aws",
			Attributes: map[string]string{"instance_type": "t3.micro"}},
	}
	fetcher := &stubFetcher{
		attrs: map[string]map[string]string{
			"aws_instance.web": {"instance_type": "t3.large"}, // changed
		},
	}

	detector := drift.NewDetector(fetcher, newTestLogger())
	results, err := detector.Check(context.Background(), resources)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Drifted {
		t.Error("expected drift to be detected")
	}
	if len(results[0].DriftedKeys) != 1 || results[0].DriftedKeys[0] != "instance_type" {
		t.Errorf("unexpected drifted keys: %v", results[0].DriftedKeys)
	}
}

func TestDetector_FetchError(t *testing.T) {
	resources := []state.Resource{
		{Name: "aws_instance.web", Type: "aws_instance", Provider: "aws",
			Attributes: map[string]string{"instance_type": "t3.micro"}},
	}
	fetcher := &stubFetcher{err: errors.New("api unavailable")}

	detector := drift.NewDetector(fetcher, newTestLogger())
	results, err := detector.Check(context.Background(), resources)

	if err != nil {
		t.Fatalf("unexpected hard error: %v", err)
	}
	if results[0].Error == "" {
		t.Error("expected error to be recorded in result")
	}
	if results[0].Drifted {
		t.Error("drifted should be false when fetch failed")
	}
}
