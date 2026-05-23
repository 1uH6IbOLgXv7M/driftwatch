// Package state provides types and utilities for reading and processing
// Terraform state files used by driftwatch.
//
// The primary entry point is the Reader interface, which abstracts the source
// of state data. LocalFileReader implements Reader for local .tfstate files.
//
// Helper functions FilterByProvider, FilterByType, and IndexByName allow
// callers to narrow down the resource list before passing it to drift
// detection logic.
package state
