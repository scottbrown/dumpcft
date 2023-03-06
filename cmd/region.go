package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func activeRegions(ctx context.Context, client *ec2.Client) ([]ec2types.Region, error) {
	resp, err := client.DescribeRegions(ctx, nil)
	if err != nil {
		return []ec2types.Region{}, err
	}

	return resp.Regions, nil
}
