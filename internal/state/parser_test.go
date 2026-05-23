package state

import (
	"encoding/json"
	"testing"
)

func buildStateJSON(version int, resources []ResourceState) []byte {
	s := TerraformState{
		Version:          version,
		TerraformVersion: "1.5.0",
		Resources:        resources,
	}
	b, _ := json.Marshal(s)
	return b
}

func TestParse_Valid(t *testing.T) {
	data := buildStateJSON(4, []ResourceState{
		{Mode: "managed", Type: "aws_instance", Name: "web", Provider: "provider[\"registry.terraform.io/hashicorp/aws\"]"},
	})

	state, err := Parse(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(state.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(state.Resources))
	}
}

func TestParse_Empty(t *testing.T) {
	_, err := Parse([]byte{})
	if err == nil {
		t.Fatal("expected error for empty data, got nil")
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte(`{not valid json}`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParse_UnsupportedVersion(t *testing.T) {
	data := buildStateJSON(3, nil)
	_, err := Parse(data)
	if err == nil {
		t.Fatal("expected error for unsupported version, got nil")
	}
}

func TestFlatResources_FiltersManaged(t *testing.T) {
	data := buildStateJSON(4, []ResourceState{
		{Mode: "managed", Type: "aws_instance", Name: "web", Provider: "aws"},
		{Mode: "data", Type: "aws_ami", Name: "latest", Provider: "aws"},
		{Mode: "managed", Type: "aws_s3_bucket", Name: "assets", Provider: "aws"},
	})

	state, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	managed := state.FlatResources()
	if len(managed) != 2 {
		t.Errorf("expected 2 managed resources, got %d", len(managed))
	}
	for _, r := range managed {
		if r.Mode != "managed" {
			t.Errorf("expected mode 'managed', got '%s'", r.Mode)
		}
	}
}

func TestFlatResources_Empty(t *testing.T) {
	data := buildStateJSON(4, []ResourceState{})
	state, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if got := state.FlatResources(); len(got) != 0 {
		t.Errorf("expected 0 resources, got %d", len(got))
	}
}
