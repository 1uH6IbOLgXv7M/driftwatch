package state_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/yourorg/driftwatch/internal/state"
)

func writeTempState(t *testing.T, s state.State) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.tfstate")
	if err != nil {
		t.Fatalf("create temp state: %v", err)
	}
	if err := json.NewEncoder(f).Encode(s); err != nil {
		t.Fatalf("encode state: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLocalFileReader_Read_Valid(t *testing.T) {
	expected := state.State{
		Version: 4,
		Resources: []state.Resource{
			{Type: "aws_instance", Name: "web", Provider: "aws", Attributes: map[string]interface{}{"id": "i-123"}},
		},
	}
	path := writeTempState(t, expected)

	reader := state.NewLocalFileReader(path)
	got, err := reader.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Version != expected.Version {
		t.Errorf("version: got %d, want %d", got.Version, expected.Version)
	}
	if len(got.Resources) != 1 {
		t.Fatalf("resources: got %d, want 1", len(got.Resources))
	}
	if got.Resources[0].Name != "web" {
		t.Errorf("resource name: got %q, want %q", got.Resources[0].Name, "web")
	}
}

func TestLocalFileReader_Read_Missing(t *testing.T) {
	reader := state.NewLocalFileReader("/nonexistent/path.tfstate")
	_, err := reader.Read()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLocalFileReader_Read_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "*.tfstate")
	f.WriteString("not json")
	f.Close()

	reader := state.NewLocalFileReader(f.Name())
	_, err := reader.Read()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLocalFileReader_Read_EmptyResources(t *testing.T) {
	// Verify that a valid state file with no resources is read without error
	// and returns an empty resource list rather than nil.
	empty := state.State{
		Version:   4,
		Resources: []state.Resource{},
	}
	path := writeTempState(t, empty)

	reader := state.NewLocalFileReader(path)
	got, err := reader.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Resources) != 0 {
		t.Errorf("resources: got %d, want 0", len(got.Resources))
	}
}
