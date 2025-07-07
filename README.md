# Formance Cloud Terraform Provider

The Formance Cloud Terraform provider allows you to manage your Formance Cloud resources via Infrastructure as Code (IaC). This provider supports managing stacks and stack modules.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Quick Start Guide](#quick-start-guide)
- [Authentication](#authentication)
- [Available Resources](#available-resources)
- [Data Sources](#data-sources)
- [Examples](#examples)
- [Full Documentation](#full-documentation)
- [Support](#support)

## Installation

### Terraform 0.13+

```hcl
terraform {
  required_providers {
    cloud = {
      source  = "formancehq/cloud"
    }
  }
}

provider "cloud" {
  # Configuration...
}
```

## Configuration

The provider can be configured in two ways:

### 1. Direct Configuration

```hcl
provider "cloud" {
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
}
```

### 2. Environment Variables

```bash
export FORMANCE_CLOUD_CLIENT_ID="your-client-id"
export FORMANCE_CLOUD_CLIENT_SECRET="your-client-secret"
```

## Quick Start Guide

Here's a minimal example to get started with the Formance Cloud provider:

```hcl
# Provider configuration
provider "cloud" {
  # Credentials can be set via environment variables
}

# Get current organization
data "cloud_current_organization" "current" {}

# Get available regions
data "cloud_regions" "default" {
  organization_id = data.cloud_current_organization.current.id
}

# Create a stack
resource "cloud_stack" "production" {
  name            = "production"
  region_id       = data.cloud_regions.default.regions[0].id
}

# Enable the ledger module
resource "cloud_stack_module" "ledger" {
  name            = "ledger"
  stack_id        = cloud_stack.production.id
}
```

## Authentication

### Getting Your Credentials

The provider uses OAuth2 authentication with client credentials. To obtain your credentials:

1. Log in to your Formance Cloud account
2. Navigate to your organization settings
3. Create a new OAuth2 application
4. Note the `client_id` and `client_secret`

### Security Best Practices

- **Never commit your credentials** in your code
- Use environment variables or a secrets manager
- Limit your credentials' permissions to the minimum required
- Rotate your secrets regularly

## Available Resources

### Stacks
- `cloud_stack` - Manages an isolated environment for your Formance services

### Modules
- `cloud_stack_module` - Enables/disables modules on a stack

## Data Sources

- `cloud_organizations` - Retrieves organization information
- `cloud_stacks` - Retrieves stack information
- `cloud_regions` - Retrieves region information
- `cloud_region_versions` - Lists available versions in a region

## Examples

### Multi-Environment Deployment

```hcl
# Get current organization
data "cloud_current_organization" "current" {}

# Get available regions
data "cloud_regions" "default" {}

# Variables for environments
variable "environments" {
  default = ["development", "staging", "production"]
}

# Create a stack for each environment
resource "cloud_stack" "env" {
  for_each        = toset(var.environments)
  name            = each.value
  region_id       = data.cloud_regions.default.regions[0].id
}

# Enable necessary modules for each stack
resource "cloud_stack_module" "ledger" {
  for_each        = cloud_stack.env
  name            = "ledger"
  stack_id        = each.value.id
}
```

## Full Documentation

For more detailed information about each resource and data source:

- [Resources Documentation](./docs/resources/)
- [Data Sources Documentation](./docs/data-sources/)
- [Complete Examples](./examples/)

## Available Modules

The following modules can be enabled on your stacks:

- **ledger** - Core accounting engine
- **payments** - Payment management and orchestration
- **webhooks** - Webhook management and distribution
- **wallets** - Digital wallet functionality
- **search** - Full-text search capabilities
- **reconciliation** - Transaction reconciliation
- **orchestration** - Workflow orchestration
- **auth** - Authentication and authorization
- **stargate** - API Gateway

## Troubleshooting

### Common Errors

#### Authentication Error
```
Error: Failed to authenticate with Formance Cloud API
```
**Solution**: Check your `client_id` and `client_secret`. Ensure they are correctly configured.

#### Permission Error
```
Error: Insufficient permissions to perform this action
```
**Solution**: Verify that your credentials have the necessary permissions for the requested action.

#### Stack Cannot Be Deleted
```
Error: Stack cannot be deleted as it contains data
```
**Solution**: Use `force_destroy = true` with caution to force deletion.

## Support

- **Issues GitHub**: [github.com/formancehq/terraform-provider-cloud/issues](https://github.com/formancehq/terraform-provider-cloud/issues)
- **API Documentation**: [docs.formance.com](https://docs.formance.com)
- **Contact**: support@formance.com

## Contributing

Contributions are welcome! See our [contribution guide](CONTRIBUTING.md) for more information.

## License

This provider is distributed under the Apache 2.0 License. See [LICENSE](LICENSE) for more details.