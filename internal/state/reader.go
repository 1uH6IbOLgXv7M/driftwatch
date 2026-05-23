package state

import (
	"encoding/json"
	"fmt"
	"os"
)

// Resource represents a single resource entry from Terraform state.
type Resource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	Attributes map[string]interface{} `json:"attributes"`
}

// State holds the parsed contents of a Terraform state file.
type State struct {
	Version   int        `json:"version"`
	Resources []Resource `json:"resources"`
}

// Reader defines the interface for reading Terraform state.
type Reader interface {
	Read() (*State, error)
}

// LocalFileReader reads state from a local .tfstate file.
type LocalFileReader struct {
	Path string
}

// NewLocalFileReader constructs a LocalFileReader for the given path.
func NewLocalFileReader(path string) *LocalFileReader {
	return &LocalFileReader{Path: path}
}

// Read opens and parses the Terraform state file.
func (r *LocalFileReader) Read() (*State, error) {
	f, err := os.Open(r.Path)
	if err != nil {
		return nil, fmt.Errorf("state: open %q: %w", r.Path, err)
	}
	defer f.Close()

	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("state: decode %q: %w", r.Path, err)
	}
	return &s, nil
}
