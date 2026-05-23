package drift_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/state"
)

func newTestReporter(t *testing.T, format string) (*drift.Reporter, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	r, err := drift.NewReporter(format, &buf)
	if err != nil {
		t.Fatalf("NewReporter(%q) unexpected error: %v", format, err)
	}
	return r, &buf
}

func sampleDriftedResources() []state.Resource {
	return []state.Resource{
		{
			Name:     "web_server",
			Type:     "aws_instance",
			Provider: "aws",
			Attributes: map[string]interface{}{
				"instance_type": "t3.large", // drifted from t3.micro
				"ami":           "ami-0abcdef1234567890",
			},
		},
		{
			Name:     "db_instance",
			Type:     "aws_db_instance",
			Provider: "aws",
			Attributes: map[string]interface{}{
				"engine_version": "14.2", // drifted from 14.1
			},
		},
	}
}

func TestReporter_Text_NoDrift(t *testing.T) {
	r, buf := newTestReporter(t, "text")
	if err := r.Report(nil); err != nil {
		t.Fatalf("Report(nil) unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected 'No drift detected' in output, got: %q", out)
	}
}

func TestReporter_Text_WithDrift(t *testing.T) {
	r, buf := newTestReporter(t, "text")
	resources := sampleDriftedResources()
	if err := r.Report(resources); err != nil {
		t.Fatalf("Report() unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "web_server") {
		t.Errorf("expected resource name 'web_server' in output, got: %q", out)
	}
	if !strings.Contains(out, "db_instance") {
		t.Errorf("expected resource name 'db_instance' in output, got: %q", out)
	}
	if !strings.Contains(out, "aws_instance") {
		t.Errorf("expected resource type 'aws_instance' in output, got: %q", out)
	}
}

func TestReporter_JSON_NoDrift(t *testing.T) {
	r, buf := newTestReporter(t, "json")
	if err := r.Report(nil); err != nil {
		t.Fatalf("Report(nil) unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	count, ok := result["drift_count"]
	if !ok {
		t.Fatal("expected 'drift_count' key in JSON output")
	}
	if count.(float64) != 0 {
		t.Errorf("expected drift_count=0, got %v", count)
	}
}

func TestReporter_JSON_WithDrift(t *testing.T) {
	r, buf := newTestReporter(t, "json")
	resources := sampleDriftedResources()
	if err := r.Report(resources); err != nil {
		t.Fatalf("Report() unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	count, ok := result["drift_count"]
	if !ok {
		t.Fatal("expected 'drift_count' key in JSON output")
	}
	if count.(float64) != float64(len(resources)) {
		t.Errorf("expected drift_count=%d, got %v", len(resources), count)
	}
	resArr, ok := result["drifted_resources"]
	if !ok {
		t.Fatal("expected 'drifted_resources' key in JSON output")
	}
	if len(resArr.([]interface{})) != len(resources) {
		t.Errorf("expected %d drifted resources, got %d", len(resources), len(resArr.([]interface{})))
	}
}

func TestReporter_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	_, err := drift.NewReporter("xml", &buf)
	if err == nil {
		t.Fatal("expected error for unsupported format 'xml', got nil")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected error to mention 'unsupported', got: %v", err)
	}
}
