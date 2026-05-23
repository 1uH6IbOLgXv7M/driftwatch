// Package drift provides functionality for comparing Terraform state
// resources against live cloud provider resources to detect configuration drift.
package drift

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/your-org/driftwatch/internal/state"
)

// Result represents the outcome of a drift check for a single resource.
type Result struct {
	// ResourceName is the Terraform resource name (e.g. "aws_instance.web").
	ResourceName string `json:"resource_name"`
	// ResourceType is the Terraform resource type (e.g. "aws_instance").
	ResourceType string `json:"resource_type"`
	// Drifted indicates whether the live state differs from the desired state.
	Drifted bool `json:"drifted"`
	// DriftedKeys lists the attribute keys whose values differ.
	DriftedKeys []string `json:"drifted_keys,omitempty"`
	// Error holds any error encountered while checking this resource.
	Error string `json:"error,omitempty"`
}

// LiveFetcher fetches live attribute values for a resource from the cloud provider.
type LiveFetcher interface {
	// Fetch returns a map of attribute key→value for the given resource.
	Fetch(ctx context.Context, r state.Resource) (map[string]string, error)
}

// Detector compares Terraform state resources against live cloud state.
type Detector struct {
	fetcher LiveFetcher
	logger  *slog.Logger
}

// NewDetector creates a new Detector using the provided LiveFetcher and logger.
func NewDetector(fetcher LiveFetcher, logger *slog.Logger) *Detector {
	return &Detector{
		fetcher: fetcher,
		logger:  logger,
	}
}

// Check evaluates each resource in the provided slice and returns drift results.
func (d *Detector) Check(ctx context.Context, resources []state.Resource) ([]Result, error) {
	results := make([]Result, 0, len(resources))

	for _, r := range resources {
		res := Result{
			ResourceName: r.Name,
			ResourceType: r.Type,
		}

		live, err := d.fetcher.Fetch(ctx, r)
		if err != nil {
			d.logger.Warn("failed to fetch live state",
				"resource", r.Name,
				"error", err,
			)
			res.Error = fmt.Sprintf("fetch error: %v", err)
			results = append(results, res)
			continue
		}

		drifted := compareAttributes(r.Attributes, live)
		res.Drifted = len(drifted) > 0
		res.DriftedKeys = drifted

		d.logger.Debug("checked resource",
			"resource", r.Name,
			"drifted", res.Drifted,
			"drifted_keys", drifted,
		)

		results = append(results, res)
	}

	return results, nil
}

// compareAttributes returns the keys whose values differ between desired and live.
func compareAttributes(desired, live map[string]string) []string {
	var drifted []string
	for k, want := range desired {
		if got, ok := live[k]; !ok || got != want {
			drifted = append(drifted, k)
		}
	}
	return drifted
}
