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
    formancecloud = {
      source  = "formancehq/formancecloud"
    }
  }
}

provider "formancecloud" {
  # Configuration...
}
```

## Configuration

The provider can be configured in two ways:

### 1. Direct Configuration

```hcl
provider "formancecloud" {
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
provider "formancecloud" {
  # Credentials can be set via environment variables
}

# Get current organization
data "formancecloud_current_organization" "current" {}

# Get available regions
data "formancecloud_regions" "default" {
  organization_id = data.formancecloud_current_organization.current.id
}

# Create a stack
resource "formancecloud_stack" "production" {
  name            = "production"
  organization_id = data.formancecloud_current_organization.current.id
  region_id       = data.formancecloud_regions.default.regions[0].id
}

# Enable the ledger module
resource "formancecloud_stack_module" "ledger" {
  name            = "ledger"
  stack_id        = formancecloud_stack.production.id
  organization_id = data.formancecloud_current_organization.current.id
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
- `formancecloud_stack` - Manages an isolated environment for your Formance services

### Modules
- `formancecloud_stack_module` - Enables/disables modules on a stack

## Data Sources

- `formancecloud_organizations` - Retrieves organization information
- `formancecloud_stacks` - Retrieves stack information
- `formancecloud_regions` - Retrieves region information
- `formancecloud_region_versions` - Lists available versions in a region

## Examples

### Multi-Environment Deployment

```hcl
# Get current organization
data "formancecloud_current_organization" "current" {}

# Get available regions
data "formancecloud_regions" "default" {
  organization_id = data.formancecloud_current_organization.current.id
}

# Variables for environments
variable "environments" {
  default = ["development", "staging", "production"]
}

# Create a stack for each environment
resource "formancecloud_stack" "env" {
  for_each        = toset(var.environments)
  name            = each.value
  organization_id = data.formancecloud_current_organization.current.id
  region_id       = data.formancecloud_regions.default.regions[0].id
}

# Enable necessary modules for each stack
resource "formancecloud_stack_module" "ledger" {
  for_each        = formancecloud_stack.env
  name            = "ledger"
  stack_id        = each.value.id
  organization_id = data.formancecloud_current_organization.current.id
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