package main

import (
	"context"
	"fmt"
	"os"

	"github.com/scottbrown/dumpcft"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dumpcft",
	Short: "Writes the templates of all CloudFormation stacks to disk",
	Long:  "You can dump the templates of all CloudFormation stacks, from any or all regions, to disk.",
	RunE:  handleRoot,
}

func validateOutputDir() error {
	_, err := os.Stat(OutputDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist.  Create it first.", OutputDir)
	}
	if err != nil {
		return fmt.Errorf("Error checking state of %s.  Cannot proceed.", OutputDir)
	}

	return nil
}

func handleRoot(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	if err := validateOutputDir(); err != nil {
		return err
	}
	fmt.Printf("Writing CloudFormation templates to directory: %s\n", OutputDir)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	ec2Client := ec2.NewFromConfig(cfg)

	regions, err := activeRegions(ctx, ec2Client)
	if err != nil {
		return err
	}

	if Regions != "" {
		regions = selectedRegions(regions)
	}

	for _, region := range regions {
		regionalCfg := cfg.Copy()
		regionalCfg.Region = *region.RegionName

		dumper := dumpcft.Dumper{
			CloudFormationClient: cloudformation.NewFromConfig(regionalCfg),
			STSClient:            sts.NewFromConfig(cfg),
			OutputDir:            OutputDir,
		}

		num, err := dumper.Dump(ctx)
		if err != nil {
			return err
		}

		fmt.Printf("%s: Wrote %d template(s)\n", *region.RegionName, num)
	}

	return nil
}
