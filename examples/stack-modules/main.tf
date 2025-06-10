terraform {
  required_providers {
    formancecloud = {
      source = "formancehq/cloud"
    }
  }
}

provider "formancecloud" {}

# TF_VAR_import_organization_id
variable "import_organization_id" {
  type = string
}

import {
  to = formancecloud_organization.default
  id = var.import_organization_id
}
resource "formancecloud_organization" "default" {
  name = "default"
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
