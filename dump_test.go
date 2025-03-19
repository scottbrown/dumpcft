package dumpcft

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Define mock implementations
type mockCloudFormationClient struct {
	describeStacksFunc func(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
	getTemplateFunc    func(ctx context.Context, params *cloudformation.GetTemplateInput, optFns ...func(*cloudformation.Options)) (*cloudformation.GetTemplateOutput, error)
}

func (m *mockCloudFormationClient) DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	return m.describeStacksFunc(ctx, params, optFns...)
}

func (m *mockCloudFormationClient) GetTemplate(ctx context.Context, params *cloudformation.GetTemplateInput, optFns ...func(*cloudformation.Options)) (*cloudformation.GetTemplateOutput, error) {
	return m.getTemplateFunc(ctx, params, optFns...)
}

type mockSTSClient struct {
	getCallerIdentityFunc func(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

func (m *mockSTSClient) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return m.getCallerIdentityFunc(ctx, params, optFns...)
}

func TestDump(t *testing.T) {
	tests := []struct {
		name              string
		outputDir         string
		describeStacksOut *cloudformation.DescribeStacksOutput
		describeStacksErr error
		getTemplateOut    *cloudformation.GetTemplateOutput
		getTemplateErr    error
		getIdentityOut    *sts.GetCallerIdentityOutput
		getIdentityErr    error
		expectError       bool
		expectedFiles     []string
		expectedCount     int
	}{
		{
			name:      "Successfully dumps a single JSON template",
			outputDir: "test-output",
			describeStacksOut: &cloudformation.DescribeStacksOutput{
				Stacks: []cfntypes.Stack{
					{
						StackName: stringPtr("test-stack"),
						StackId:   stringPtr("arn:aws:cloudformation:us-west-2:123456789012:stack/test-stack/abcdef"),
					},
				},
			},
			getTemplateOut: &cloudformation.GetTemplateOutput{
				TemplateBody: stringPtr(`{"Resources": {"TestResource": {"Type": "AWS::S3::Bucket"}}}`),
			},
			getIdentityOut: &sts.GetCallerIdentityOutput{
				Account: stringPtr("123456789012"),
			},
			expectedFiles: []string{
				"123456789012.us-west-2.test-stack.cfn.json",
			},
			expectedCount: 1,
		},
		{
			name:      "Successfully dumps a single YAML template",
			outputDir: "test-output",
			describeStacksOut: &cloudformation.DescribeStacksOutput{
				Stacks: []cfntypes.Stack{
					{
						StackName: stringPtr("test-stack"),
						StackId:   stringPtr("arn:aws:cloudformation:us-west-2:123456789012:stack/test-stack/abcdef"),
					},
				},
			},
			getTemplateOut: &cloudformation.GetTemplateOutput{
				TemplateBody: stringPtr("Resources:\n  TestResource:\n    Type: AWS::S3::Bucket"),
			},
			getIdentityOut: &sts.GetCallerIdentityOutput{
				Account: stringPtr("123456789012"),
			},
			expectedFiles: []string{
				"123456789012.us-west-2.test-stack.cfn.yaml",
			},
			expectedCount: 1,
		},
		{
			name:              "Error in DescribeStacks",
			outputDir:         "test-output",
			describeStacksErr: errors.New("describe stacks error"),
			expectError:       true,
		},
		{
			name:      "Error in GetTemplate",
			outputDir: "test-output",
			describeStacksOut: &cloudformation.DescribeStacksOutput{
				Stacks: []cfntypes.Stack{
					{
						StackName: stringPtr("test-stack"),
						StackId:   stringPtr("arn:aws:cloudformation:us-west-2:123456789012:stack/test-stack/abcdef"),
					},
				},
			},
			getTemplateErr: errors.New("get template error"),
			getIdentityOut: &sts.GetCallerIdentityOutput{
				Account: stringPtr("123456789012"),
			},
			expectError: true,
		},
		{
			name:      "Error in GetCallerIdentity",
			outputDir: "test-output",
			describeStacksOut: &cloudformation.DescribeStacksOutput{
				Stacks: []cfntypes.Stack{
					{
						StackName: stringPtr("test-stack"),
						StackId:   stringPtr("arn:aws:cloudformation:us-west-2:123456789012:stack/test-stack/abcdef"),
					},
				},
			},
			getTemplateOut: &cloudformation.GetTemplateOutput{
				TemplateBody: stringPtr(`{"Resources": {"TestResource": {"Type": "AWS::S3::Bucket"}}}`),
			},
			getIdentityErr: errors.New("get caller identity error"),
			expectError:    true,
		},
		{
			name:      "Successful pagination",
			outputDir: "test-output",
			describeStacksOut: &cloudformation.DescribeStacksOutput{
				Stacks: []cfntypes.Stack{
					{
						StackName: stringPtr("test-stack-1"),
						StackId:   stringPtr("arn:aws:cloudformation:us-west-2:123456789012:stack/test-stack-1/abcdef"),
					},
				},
				NextToken: stringPtr("next-token"),
			},
			getTemplateOut: &cloudformation.GetTemplateOutput{
				TemplateBody: stringPtr(`{"Resources": {"TestResource": {"Type": "AWS::S3::Bucket"}}}`),
			},
			getIdentityOut: &sts.GetCallerIdentityOutput{
				Account: stringPtr("123456789012"),
			},
			expectedFiles: []string{
				"123456789012.us-west-2.test-stack-1.cfn.json",
				"123456789012.us-west-2.test-stack-2.cfn.json",
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test directory
			err := os.MkdirAll(tt.outputDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}
			defer os.RemoveAll(tt.outputDir)

			// Create mock CloudFormation client
			mockCfn := &mockCloudFormationClient{
				describeStacksFunc: func(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
					if tt.describeStacksErr != nil {
						return nil, tt.describeStacksErr
					}

					// Handle pagination case
					if tt.name == "Successful pagination" && params.NextToken != nil {
						// Return second page
						return &cloudformation.DescribeStacksOutput{
							Stacks: []cfntypes.Stack{
								{
									StackName: stringPtr("test-stack-2"),
									StackId:   stringPtr("arn:aws:cloudformation:us-west-2:123456789012:stack/test-stack-2/ghijkl"),
								},
							},
						}, nil
					}

					return tt.describeStacksOut, nil
				},
				getTemplateFunc: func(ctx context.Context, params *cloudformation.GetTemplateInput, optFns ...func(*cloudformation.Options)) (*cloudformation.GetTemplateOutput, error) {
					if tt.getTemplateErr != nil {
						return nil, tt.getTemplateErr
					}
					return tt.getTemplateOut, nil
				},
			}

			// Create mock STS client
			mockSTS := &mockSTSClient{
				getCallerIdentityFunc: func(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
					if tt.getIdentityErr != nil {
						return nil, tt.getIdentityErr
					}
					return tt.getIdentityOut, nil
				},
			}

			// Create dumper
			dumper := Dumper{
				CloudFormationClient: mockCfn,
				STSClient:            mockSTS,
				OutputDir:            tt.outputDir,
			}

			// Call Dump
			count, err := dumper.Dump(context.Background())

			// Check for expected errors
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got nil")
				}
				return
			}

			// Check for unexpected errors
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check count
			if count != tt.expectedCount {
				t.Errorf("Expected count %d but got %d", tt.expectedCount, count)
			}

			// Check for expected files
			for _, filename := range tt.expectedFiles {
				path := filepath.Join(tt.outputDir, filename)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("Expected file %s does not exist", path)
				}
			}
		})
	}
}

func TestIsJSON(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Valid JSON",
			content:  `{"key": "value"}`,
			expected: true,
		},
		{
			name:     "Valid JSON array",
			content:  `[1, 2, 3]`,
			expected: true,
		},
		{
			name:     "Valid JSON number",
			content:  `42`,
			expected: true,
		},
		{
			name:     "Valid JSON string",
			content:  `"string"`,
			expected: true,
		},
		{
			name:     "Valid JSON boolean",
			content:  `true`,
			expected: true,
		},
		{
			name:     "Valid JSON null",
			content:  `null`,
			expected: true,
		},
		{
			name:     "YAML content",
			content:  "key: value",
			expected: false,
		},
		{
			name: "Complex YAML",
			content: `Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-bucket`,
			expected: false,
		},
		{
			name:     "Empty string",
			content:  "",
			expected: false,
		},
		{
			name:     "Invalid JSON",
			content:  `{"key": value}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isJSON([]byte(tt.content))
			if result != tt.expected {
				t.Errorf("Expected isJSON to return %v but got %v for content: %s", tt.expected, result, tt.content)
			}
		})
	}
}

func TestBuildFilename(t *testing.T) {
	tests := []struct {
		name      string
		accountId string
		region    string
		stackName string
		ext       string
		outputDir string
		expected  string
	}{
		{
			name:      "Basic filename",
			accountId: "123456789012",
			region:    "us-west-2",
			stackName: "test-stack",
			ext:       "json",
			outputDir: "output",
			expected:  "output/123456789012.us-west-2.test-stack.cfn.json",
		},
		{
			name:      "Stack name with special characters",
			accountId: "123456789012",
			region:    "us-east-1",
			stackName: "stack-with_special.chars",
			ext:       "yaml",
			outputDir: "/tmp/output",
			expected:  "/tmp/output/123456789012.us-east-1.stack-with_special.chars.cfn.yaml",
		},
		{
			name:      "Different directory",
			accountId: "987654321098",
			region:    "eu-west-1",
			stackName: "production-stack",
			ext:       "json",
			outputDir: "./templates",
			expected:  "./templates/987654321098.eu-west-1.production-stack.cfn.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dumper := Dumper{
				OutputDir: tt.outputDir,
			}
			result := dumper.buildFilename(tt.accountId, tt.region, tt.stackName, tt.ext)
			if result != tt.expected {
				t.Errorf("Expected filename %q but got %q", tt.expected, result)
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
