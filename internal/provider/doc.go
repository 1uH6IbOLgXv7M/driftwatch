// Package provider defines the Fetcher interface and supporting types
// used by the drift detector to retrieve live resource attributes from
// cloud providers.
//
// Implementations of Fetcher are provider-specific (e.g. AWS, GCP) and
// are registered at startup based on the active configuration. A
// MockFetcher is provided for use in unit tests.
package provider
