package main

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"golang.org/x/exp/slices"
)

func activeRegions(ctx context.Context, client *ec2.Client) ([]types.Region, error) {
	resp, err := client.DescribeRegions(ctx, nil)
	if err != nil {
		return []types.Region{}, err
	}

	return resp.Regions, nil
}

func selectedRegions(activeRegions []types.Region) []types.Region {
	// parse regions flag into parts
	targetRegions := strings.Split(Regions, ",")

	// remove the regions that the user didn't target
	var selectedRegions []types.Region
	for _, r := range activeRegions {
		if slices.Contains(targetRegions, *r.RegionName) {
			selectedRegions = append(selectedRegions, r)
		}
	}

	return selectedRegions
}
