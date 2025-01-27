package dumpcft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"gopkg.in/yaml.v3"
)

type Dumper struct {
	CloudFormationClient *cloudformation.Client
	STSClient            *sts.Client
	OutputDir            string
}

func (d Dumper) Dump(ctx context.Context) (int, error) {
	var stacks []cfntypes.Stack
	var nextToken *string
	for {
		input := cloudformation.DescribeStacksInput{
			NextToken: nextToken,
		}
		resp, err := d.CloudFormationClient.DescribeStacks(ctx, &input)

		if err != nil {
			return 0, err
		}

		for _, s := range resp.Stacks {
			stacks = append(stacks, s)
		}

		if nextToken == nil {
			break
		}
	}

	for _, stack := range stacks {
		input := cloudformation.GetTemplateInput{
			StackName: stack.StackName,
		}
		resp, err := d.CloudFormationClient.GetTemplate(ctx, &input)
		if err != nil {
			return 0, err
		}

		arn, err := arn.Parse(*stack.StackId)
		if err != nil {
			return 0, err
		}

		accountId, err := d.awsAccountId(ctx)
		if err != nil {
			return 0, err
		}

		templateBody := *resp.TemplateBody

		var filename string
		var formattedTemplate string

		if isJSON([]byte(templateBody)) {
			filename = d.buildFilename(accountId, arn.Region, *stack.StackName, "json")
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(templateBody), "", "  "); err != nil {
				return 0, fmt.Errorf("failed to format JSON: %w", err)
			}
			formattedTemplate = prettyJSON.String()
		} else {
			filename = d.buildFilename(accountId, arn.Region, *stack.StackName, "yaml")
			var v interface{}
			if err := yaml.Unmarshal([]byte(templateBody), &v); err != nil {
				return 0, fmt.Errorf("failed to parse YAML: %w", err)
			}
			yamlBytes, err := yaml.Marshal(v)
			if err != nil {
				return 0, fmt.Errorf("failed to format YAML: %w", err)
			}
			formattedTemplate = string(yamlBytes)
		}

		file, err := os.Create(filename) // #nosec G304 -- AWS is in trust boundary
		if err != nil {
			return 0, err
		}
		defer func() {
			if err := file.Close(); err != nil {
				// TODO print error
			}
		}()

		_, err = file.WriteString(formattedTemplate)
		if err != nil {
			return 0, err
		}
	}

	return len(stacks), nil
}

func (d Dumper) buildFilename(accountId, region, stackName, ext string) string {
	return fmt.Sprintf("%s/%s.%s.%s.cfn.%s", d.OutputDir, accountId, region, stackName, ext)
}

func (d Dumper) awsAccountId(ctx context.Context) (string, error) {
	resp, err := d.STSClient.GetCallerIdentity(ctx, nil)
	if err != nil {
		return "", err
	}

	return *resp.Account, nil
}

func isJSON(content []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(content, &js) == nil
}
