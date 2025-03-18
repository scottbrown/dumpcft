# dumpcft

`dumpcft` is a streamlined CLI tool for AWS CloudFormation template extraction. It allows you to export all CloudFormation templates from your AWS account to local files for documentation, backup, or analysis purposes.

## Features

- Export CloudFormation templates from all or selected AWS regions
- Automatic file format detection (YAML or JSON)
- Consistent file naming with account ID, region, and stack name
- Supports multiple AWS profiles

## Installation

### Pre-Built Binaries

Download the latest release from the [releases page](https://github.com/scottbrown/dumpcft/releases) for your platform:

- Linux (amd64, arm64)
- MacOS (amd64, arm64)
- Windows (amd64, arm64)

### From Source

If you have Go 1.24+ installed:

```bash
go install github.com/scottbrown/dumpcft/cmd@latest
```

## Usage

### Basic Usage

Create an output directory and run `dumpcft`:

```bash
mkdir templates
dumpcft --profile PROFILE templates
```

This will export all CloudFormation templates from all active regions to the `templates` directory.

### Command Line Options

```
Usage:
  dumpcft [flags] [output-directory]

Flags:
  -h, --help                  help for dumpcft
  -o, --output-dir string     The directory where templates are persisted to disk (default "templates")
  -p, --profile string        The AWS profile to use
  -r, --regions string        One or more comma-delimited regions to dump
  -v, --verbose               Shows debug output
```

### Examples

### Specify Output Directory

```bash
dumpcft --profile your-aws-profile -o /tmp/cfn-templates
```

### Limit to Specific Regions

```bash
dumpcft --profile your-aws-profile --regions us-east-1,us-west-2 templates
```

## File Format

Templates are saved with the following naming convention:

```
<AWS_ACCOUNT_ID>.<REGION>.<STACK_NAME>.cfn.<FORMAT>
```

For example:

- `012345678901.us-east-1.my-application-stack.cfn.yaml`
- `012345678901.us-west-2.database-stack.cfn.json`

The format extension (`.yaml` or `.json`) is automatically determined by analyzing the template content.

## Use Cases

- **Infrastructure Auditing**: Snapshot all your CloudFormation templates for compliance reviews
- **Documentation**: Export templates to include in architecture documentation
- **Code Archaeology**: Analyze historical infrastructure changes by comparing template snapshots
- **Backup**: Store copies of critical infrastructure templates for disaster recovery
- **Migration**: Prepare for cross-account or cross-region migrations

## AWS Permissions

The following AWS permissions are required:

- cloudformation:DescribeStacks
- cloudformation:GetTemplate
- ec2:DescribeRegions
- sts:GetCallerIdentity

## Development

### Prerequisites

- Go 1.24 or higher
- AWS account for integration testing

### Build

```bash
make build
```
### Test

```bash
make test
```

### Security Checks

```bash
make check
```

### Format Code

```bash
make fmt
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
1. Create a feature branch (`git checkout -b feature/amazing-feature`)
1. Commit your changes (`git commit -m 'Add amazing feature'`)
1. Push to your branch `(git push origin feature/amazing-feature`)
1. Open a Pull Request

Before submitting:

- Run `make fmt` to ensure consistent code formatting
- Run `make test` to verify all tests pass
- Run `make check` to perform security and code quality checks

## License

This project is licensed under the MIT License – see the [LICENSE](LICENSE) file for details.
