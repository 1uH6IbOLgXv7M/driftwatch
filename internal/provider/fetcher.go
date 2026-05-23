// Package provider defines interfaces and helpers for fetching live
// cloud resource attributes from various infrastructure providers.
package provider

import "context"

// Attributes is a map of resource attribute key/value pairs as observed
// in the live cloud environment.
type Attributes map[string]interface{}

// ResourceRef uniquely identifies a cloud resource by its type and name.
type ResourceRef struct {
	Type string
	Name string
}

// Fetcher retrieves live attributes for a given resource from a cloud
// provider. Implementations must be safe for concurrent use.
type Fetcher interface {
	// Fetch returns the live attributes for the given resource reference.
	// It returns an error if the resource cannot be found or the provider
	// call fails.
	Fetch(ctx context.Context, ref ResourceRef) (Attributes, error)
}

// ErrNotFound is returned by a Fetcher when the requested resource does
// not exist in the live environment.
type ErrNotFound struct {
	Ref ResourceRef
}

func (e *ErrNotFound) Error() string {
	return "provider: resource not found: " + e.Ref.Type + "/" + e.Ref.Name
}
