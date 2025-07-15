# Complete example of Formance Cloud configuration with Terraform
# This example shows how to configure a complete infrastructure

terraform {
  required_providers {
    cloud = {
      source  = "formancehq/cloud"
      version = "~> 1.0"
    }
  }
}

# Provider configuration (credentials via environment variables)
provider "cloud" {}

# Configuration variables
variable "organization_name" {
  description = "Organization name"
  type        = string
  default     = "my-company"
}

variable "domain" {
  description = "Organization domain"
  type        = string
  default     = "mycompany.com"
}

variable "environments" {
  description = "List of environments to create"
  type        = list(string)
  default     = ["development", "staging", "production"]
}

variable "team_members" {
  description = "Team members and their roles"
  type = map(object({
    email = string
    role  = string
  }))
  default = {
    lead_dev = {
      email = "lead@mycompany.com"
      role  = "WRITE"
    }
    dev1 = {
      email = "dev1@mycompany.com"
      role  = "WRITE"
    }
    dev2 = {
      email = "dev2@mycompany.com"
      role  = "WRITE"
    }
    analyst = {
      email = "analyst@mycompany.com"
      role  = "READ"
    }
  }
}

variable "modules_to_enable" {
  description = "Modules to enable on stacks"
  type        = list(string)
  default = [
    "ledger",
    "payments",
    "webhooks",
    "wallets",
    "auth",
    "stargate"
  ]
}

# Organization creation
resource "cloud_organization" "main" {
  name                        = var.organization_name
  domain                      = var.domain
  default_organization_access = "READ"
  default_stack_access        = "NONE" # Explicit access required
}

# Creation of a private region for Europe
resource "cloud_region" "europe" {
  name = "europe-west"
}

# Creation of a private region for US (optional)
resource "cloud_region" "us" {
  name = "us-east"
}

# Fetching available versions
data "cloud_region_versions" "europe" {
  id = cloud_region.europe.id
}

# Creating stacks for each environment
resource "cloud_stack" "environments" {
  for_each = toset(var.environments)

  name      = each.value
  region_id = cloud_region.europe.id

  # Use latest stable version for dev/staging, fixed version for prod
  version = each.value == "production" ? "v2.0.0" : data.cloud_region_versions.europe.versions[0].name

  # Protection against accidental deletion in production
  force_destroy = each.value != "production"

  lifecycle {
    # Prevent accidental deletion of production stack
    prevent_destroy = false # Set to true in real production
  }
}

# Enabling modules on each stack
resource "cloud_stack_module" "modules" {
  for_each = {
    for pair in setproduct(keys(cloud_stack.environments), var.modules_to_enable) :
    "${pair[0]}-${pair[1]}" => {
      stack_key = pair[0]
      module    = pair[1]
    }
  }

  name     = each.value.module
  stack_id = cloud_stack.environments[each.value.stack_key].id

  # Modules have dependencies, ensure they are created in the right order
  depends_on = [
    cloud_stack.environments
  ]
}

# Adding members to the organization
resource "cloud_organization_member" "team" {
  for_each = var.team_members

  email = each.value.email
  role  = each.value.role
}

# Stack access configuration
locals {
  # Access matrix: who has access to which environment
  stack_access = {
    # Everyone has access to dev
    development = {
      for name, member in var.team_members : name => member.role
    }
    # Only devs have access to staging
    staging = {
      for name, member in var.team_members : name => member.role
      if member.role == "WRITE"
    }
    # Restricted access to production
    production = {
      lead_dev = "WRITE"
      analyst  = "READ"
    }
  }
}

# Assigning access to stacks
resource "cloud_stack_member" "access" {
  for_each = {
    for item in flatten([
      for env, access in local.stack_access : [
        for member_name, role in access : {
          key         = "${env}-${member_name}"
          env         = env
          member_name = member_name
          role        = role
        }
      ]
    ]) : item.key => item
  }

  stack_id = cloud_stack.environments[each.value.env].id
  user_id  = cloud_organization_member.team[each.value.member_name].user_id
  role     = each.value.role
}

# Dedicated stack for integration tests (CI/CD)
resource "cloud_stack" "ci" {
  name          = "ci-testing"
  region_id     = cloud_region.europe.id
  force_destroy = true # Can be deleted without confirmation
}

# Minimal modules for CI tests
resource "cloud_stack_module" "ci_modules" {
  for_each = toset(["ledger", "auth"])

  name     = each.value
  stack_id = cloud_stack.ci.id
}

# Useful outputs
output "organization_id" {
  description = "ID of the created organization"
  value       = cloud_organization.main.id
}

output "stack_urls" {
  description = "URLs of created stacks"
  value = {
    for name, stack in cloud_stack.environments : name => stack.uri
  }
}

output "region_endpoints" {
  description = "Region endpoints"
  value = {
    europe = cloud_region.europe.base_url
    us     = cloud_region.us.base_url
  }
}

output "ci_stack_url" {
  description = "CI stack URL for automated tests"
  value       = cloud_stack.ci.uri
}

# Important note about region secret (displayed only at creation)
output "region_secret_note" {
  description = "Note about region secrets"
  value       = "Region secrets are available only at creation time. Store them securely."
}