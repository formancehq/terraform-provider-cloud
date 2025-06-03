terraform {
  required_providers {
    formancecloud = {
      source = "formancehq/cloud"
    }
  }
}

provider "formancecloud" {}

resource "formancecloud_organization" "default" {
  name = "formancehq"
}

variable "region_datasource_name" {
  type = string
}

output "region_datasource_name" {
  value = data.formancecloud_regions.dev.id
}

data "formancecloud_regions" "dev" {
  name            = var.region_datasource_name
  organization_id = formancecloud_organization.default.id
}

resource "formancecloud_stack" "default" {
  name            = "test-stack"
  organization_id = formancecloud_organization.default.id
  region_id       = data.formancecloud_regions.dev.id
  version         = "default"
}

resource "formancecloud_stack_module" "default_webhooks" {
  name            = "webhooks"
  organization_id = formancecloud_organization.default.id
  stack_id        = formancecloud_stack.default.id
}

resource "formancecloud_stack_module" "default_reconciliation" {
  name            = "reconciliation"
  organization_id = formancecloud_organization.default.id
  stack_id        = formancecloud_stack.default.id
}

resource "formancecloud_stack_module" "default_orchestration" {
  name            = "orchestration"
  organization_id = formancecloud_organization.default.id
  stack_id        = formancecloud_stack.default.id
}
