package provider

import (
	"context"
	"fmt"
)

// MockFetcher is an in-memory Fetcher used in tests. It maps
// "<Type>/<Name>" keys to pre-configured Attributes.
type MockFetcher struct {
	// Resources holds the attributes keyed by "<Type>/<Name>".
	Resources map[string]Attributes
	// Err, when non-nil, is returned for every Fetch call regardless of
	// the resource reference.
	Err error
}

// NewMockFetcher returns a MockFetcher pre-populated with the provided
// resource map.
func NewMockFetcher(resources map[string]Attributes) *MockFetcher {
	return &MockFetcher{Resources: resources}
}

// Fetch implements Fetcher. It returns the pre-configured attributes for
// the given ResourceRef or ErrNotFound if the key is absent.
func (m *MockFetcher) Fetch(_ context.Context, ref ResourceRef) (Attributes, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	key := fmt.Sprintf("%s/%s", ref.Type, ref.Name)
	attrs, ok := m.Resources[key]
	if !ok {
		return nil, &ErrNotFound{Ref: ref}
	}
	return attrs, nil
}
