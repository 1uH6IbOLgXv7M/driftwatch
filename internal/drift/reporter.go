// Package drift provides functionality for detecting and reporting
// configuration drift between Terraform state and live cloud resources.
package drift

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// DriftReport summarizes the results of a drift detection run.
type DriftReport struct {
	GeneratedAt  time.Time      `json:"generated_at"`
	Provider     string         `json:"provider"`
	TotalChecked int            `json:"total_checked"`
	DriftedCount int            `json:"drifted_count"`
	CleanCount   int            `json:"clean_count"`
	Results      []DriftResult  `json:"results"`
}

// DriftResult holds the drift status for a single resource.
type DriftResult struct {
	ResourceName string            `json:"resource_name"`
	ResourceType string            `json:"resource_type"`
	Drifted      bool              `json:"drifted"`
	Diffs        []AttributeDiff   `json:"diffs,omitempty"`
}

// AttributeDiff describes a single attribute mismatch between
// the Terraform state and the live cloud resource.
type AttributeDiff struct {
	Attribute string `json:"attribute"`
	Expected  string `json:"expected"`
	Actual    string `json:"actual"`
}

// ReportFormat controls the output format of a DriftReport.
type ReportFormat string

const (
	FormatText ReportFormat = "text"
	FormatJSON ReportFormat = "json"
)

// Reporter writes drift reports to an io.Writer in the configured format.
type Reporter struct {
	out    io.Writer
	format ReportFormat
}

// NewReporter creates a Reporter that writes to out using the given format.
// Defaults to text format if an unrecognised format string is provided.
func NewReporter(out io.Writer, format ReportFormat) *Reporter {
	if format != FormatJSON {
		format = FormatText
	}
	return &Reporter{out: out, format: format}
}

// Write renders the report to the configured writer.
func (r *Reporter) Write(report DriftReport) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(report)
	default:
		return r.writeText(report)
	}
}

func (r *Reporter) writeJSON(report DriftReport) error {
	enc := json.NewEncoder(r.out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("reporter: encoding JSON: %w", err)
	}
	return nil
}

func (r *Reporter) writeText(report DriftReport) error {
	fmt.Fprintf(r.out, "Drift Report — %s\n", report.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(r.out, "Provider : %s\n", report.Provider)
	fmt.Fprintf(r.out, "Checked  : %d  Drifted: %d  Clean: %d\n",
		report.TotalChecked, report.DriftedCount, report.CleanCount)
	fmt.Fprintln(r.out, strings.Repeat("-", 60))

	for _, res := range report.Results {
		status := "OK"
		if res.Drifted {
			status = "DRIFT"
		}
		fmt.Fprintf(r.out, "[%s] %s (%s)\n", status, res.ResourceName, res.ResourceType)
		for _, d := range res.Diffs {
			fmt.Fprintf(r.out, "      attribute : %s\n", d.Attribute)
			fmt.Fprintf(r.out, "      expected  : %s\n", d.Expected)
			fmt.Fprintf(r.out, "      actual    : %s\n", d.Actual)
		}
	}

	return nil
}
