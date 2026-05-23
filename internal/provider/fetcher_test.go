package provider_test

import (
	"context"
	"errors"
	"testing"

	"github.com/example/driftwatch/internal/provider"
)

func TestMockFetcher_Found(t *testing.T) {
	f := provider.NewMockFetcher(map[string]provider.Attributes{
		"aws_instance/web": {"instance_type": "t3.micro", "ami": "ami-123"},
	})

	attrs, err := f.Fetch(context.Background(), provider.ResourceRef{
		Type: "aws_instance",
		Name: "web",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["instance_type"] != "t3.micro" {
		t.Errorf("expected t3.micro, got %v", attrs["instance_type"])
	}
}

func TestMockFetcher_NotFound(t *testing.T) {
	f := provider.NewMockFetcher(map[string]provider.Attributes{})

	_, err := f.Fetch(context.Background(), provider.ResourceRef{
		Type: "aws_s3_bucket",
		Name: "logs",
	})
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	var notFound *provider.ErrNotFound
	if !errors.As(err, &notFound) {
		t.Errorf("expected *ErrNotFound, got %T", err)
	}
	if notFound.Ref.Name != "logs" {
		t.Errorf("unexpected ref name: %s", notFound.Ref.Name)
	}
}

func TestMockFetcher_ForcedError(t *testing.T) {
	sentinel := errors.New("provider unavailable")
	f := &provider.MockFetcher{Err: sentinel}

	_, err := f.Fetch(context.Background(), provider.ResourceRef{
		Type: "aws_instance",
		Name: "web",
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestErrNotFound_Error(t *testing.T) {
	e := &provider.ErrNotFound{Ref: provider.ResourceRef{Type: "aws_vpc", Name: "main"}}
	want := "provider: resource not found: aws_vpc/main"
	if e.Error() != want {
		t.Errorf("got %q, want %q", e.Error(), want)
	}
}
