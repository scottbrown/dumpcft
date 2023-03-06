# dumpcft

A simple CLI tool that will take an AWS profile, and an output directory,
and dump all of the CloudFormation templates in every region into files
on the local disk.

How is this helpful?  I don't know about any general case, but I needed it
to solve a specific problem around code archaeology.

## Usage

Dump templates from all regions into a pre-existing directory called `templates`.

```bash
mkdir templates
dumpcft --profile PROFILE templates
Writing CloudFormation templates to directory: templates.
```

Specify a set of regions in a comma-delimited list.

```bash
mkdir templates
dumpcft --profile PROFILE --regions us-east-1,us-west-1 templates
Writing CloudFormation templates to directory: templates.
```

Write the templates to a directory of your choice.

```bash
dumpcft --profile PROFILE -o /tmp/templates
Writing CloudFormation templates to directory: /tmp/templates
```

All CloudFormation templates dumped by this tool will have the filename
format:

```
AWS_ACCOUNT_ID.REGION.STACK_NAME.[yaml|json]
```

An example of which would be `01234567890.ca-central-1.example-app.yaml`
for a CloudFormation stack named `example-app` residing in `ca-central-1`
within the AWS account with ID `01234567890`.

## Contributing

1. Fork the repository.
1. Make your change.
1. Ensure `make fmt` is run.
1. Ensure `make test` completes successfully.
1. Make a Pull Request.
1. Add a clear, concise summary in the Pull Request.

## License

[MIT](LICENSE)
