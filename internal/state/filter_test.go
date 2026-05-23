package state_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/state"
)

var testResources = []state.Resource{
	{Type: "aws_instance", Name: "web", Provider: "aws"},
	{Type: "aws_s3_bucket", Name: "assets", Provider: "aws"},
	{Type: "google_compute_instance", Name: "vm", Provider: "google"},
}

func TestFilterByProvider(t *testing.T) {
	got := state.FilterByProvider(testResources, "aws")
	if len(got) != 2 {
		t.Errorf("got %d resources, want 2", len(got))
	}
}

func TestFilterByProvider_Empty(t *testing.T) {
	got := state.FilterByProvider(testResources, "")
	if len(got) != len(testResources) {
		t.Errorf("empty filter: got %d, want %d", len(got), len(testResources))
	}
}

func TestFilterByType(t *testing.T) {
	got := state.FilterByType(testResources, "aws_instance")
	if len(got) != 1 {
		t.Errorf("got %d resources, want 1", len(got))
	}
	if got[0].Name != "web" {
		t.Errorf("got name %q, want %q", got[0].Name, "web")
	}
}

func TestFilterByType_NoMatch(t *testing.T) {
	got := state.FilterByType(testResources, "azure_vm")
	if len(got) != 0 {
		t.Errorf("expected 0 results, got %d", len(got))
	}
}

func TestIndexByName(t *testing.T) {
	idx := state.IndexByName(testResources)
	if len(idx) != 3 {
		t.Errorf("index length: got %d, want 3", len(idx))
	}
	if _, ok := idx["vm"]; !ok {
		t.Error("expected key 'vm' in index")
	}
}
