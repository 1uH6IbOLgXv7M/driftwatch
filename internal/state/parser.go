package state

import (
	"encoding/json"
	"fmt"
)

// TerraformState represents the top-level structure of a Terraform state file.
type TerraformState struct {
	Version          int               `json:"version"`
	TerraformVersion string            `json:"terraform_version"`
	Resources        []ResourceState   `json:"resources"`
}

// ResourceState represents a single resource block in Terraform state.
type ResourceState struct {
	Module    string          `json:"module,omitempty"`
	Mode      string          `json:"mode"`
	Type      string          `json:"type"`
	Name      string          `json:"name"`
	Provider  string          `json:"provider"`
	Instances []InstanceState `json:"instances"`
}

// InstanceState holds the attributes of a single resource instance.
type InstanceState struct {
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
	SensitiveAttributes []interface{}   `json:"sensitive_attributes,omitempty"`
}

// Parse deserializes raw JSON bytes into a TerraformState.
// It returns an error if the bytes are not valid JSON or the version is unsupported.
func Parse(data []byte) (*TerraformState, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("state data is empty")
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse terraform state: %w", err)
	}

	if state.Version < 4 {
		return nil, fmt.Errorf("unsupported state version %d: only version 4+ is supported", state.Version)
	}

	return &state, nil
}

// FlatResources returns a flat slice of ResourceState containing only
// managed resources (mode == "managed").
func (s *TerraformState) FlatResources() []ResourceState {
	managed := make([]ResourceState, 0, len(s.Resources))
	for _, r := range s.Resources {
		if r.Mode == "managed" {
			managed = append(managed, r)
		}
	}
	return managed
}
