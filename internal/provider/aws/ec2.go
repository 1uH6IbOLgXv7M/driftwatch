// Package aws provides cloud resource fetchers for Amazon Web Services.
package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/yourorg/driftwatch/internal/provider"
)

// EC2Fetcher fetches live attributes for AWS EC2 instances.
type EC2Fetcher struct {
	client *ec2.Client
	region string
}

// NewEC2Fetcher creates an EC2Fetcher using the default AWS credential chain.
// The region is loaded from the environment or the provided AWS config profile.
func NewEC2Fetcher(ctx context.Context, region string) (*EC2Fetcher, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("aws/ec2: load config: %w", err)
	}

	return &EC2Fetcher{
		client: ec2.NewFromConfig(cfg),
		region: region,
	}, nil
}

// Fetch retrieves the live attributes of an EC2 instance identified by resourceName.
// resourceName must be the instance ID (e.g. "i-0abc123def456").
// Returns provider.ErrNotFound when no matching instance exists.
func (f *EC2Fetcher) Fetch(ctx context.Context, resourceName string) (map[string]any, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{resourceName},
	}

	out, err := f.client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("aws/ec2: describe instance %q: %w", resourceName, err)
	}

	instance, err := extractInstance(out.Reservations, resourceName)
	if err != nil {
		return nil, err
	}

	return instanceAttributes(instance), nil
}

// extractInstance finds the first instance across all reservations.
func extractInstance(reservations []types.Reservation, id string) (types.Instance, error) {
	for _, r := range reservations {
		for _, inst := range r.Instances {
			if aws.ToString(inst.InstanceId) == id {
				return inst, nil
			}
		}
	}
	return types.Instance{}, fmt.Errorf("aws/ec2: instance %q: %w", id, provider.ErrNotFound)
}

// instanceAttributes converts an EC2 instance into a flat attribute map that
// mirrors the keys used in Terraform's aws_instance resource schema.
func instanceAttributes(inst types.Instance) map[string]any {
	attrs := map[string]any{
		"instance_type": string(inst.InstanceType),
		"ami":           aws.ToString(inst.ImageId),
		"key_name":      aws.ToString(inst.KeyName),
		"subnet_id":     aws.ToString(inst.SubnetId),
		"vpc_id":        aws.ToString(inst.VpcId),
		"private_ip":    aws.ToString(inst.PrivateIpAddress),
		"public_ip":     aws.ToString(inst.PublicIpAddress),
		"state":         string(inst.State.Name),
	}

	// Flatten tags into "tags.<key>" entries for easy comparison.
	for _, tag := range inst.Tags {
		key := fmt.Sprintf("tags.%s", aws.ToString(tag.Key))
		attrs[key] = aws.ToString(tag.Value)
	}

	return attrs
}
