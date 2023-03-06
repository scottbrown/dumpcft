package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/scottbrown/dumpcft"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
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

func selectedRegions(activeRegions []ec2types.Region) []ec2types.Region {
	// parse regions flag into parts
	targetRegions := strings.Split(Regions, ",")

	// remove the regions that the user didn't target
	var selectedRegions []ec2types.Region
	for _, r := range activeRegions {
		if slices.Contains(targetRegions, *r.RegionName) {
			selectedRegions = append(selectedRegions, r)
		}
	}

	return selectedRegions
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
