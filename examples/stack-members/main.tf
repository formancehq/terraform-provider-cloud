terraform {
  required_providers {
    cloud = {
      source = "formancehq/cloud"
    }
  }
}

provider "cloud" {}

# TF_VAR_import_organization_id
variable "import_organization_id" {
  type = string
}

import {
  to = cloud_organization.default
  id = var.import_organization_id
}
resource "cloud_organization" "default" {
  name = "default"
}

resource "cloud_organization_member" "default" {
  email = "example@formance.dev"
  role  = "ADMIN"
}

resource "cloud_region" "dev" {}

resource "cloud_stack" "default" {
  name      = "default"
  region_id = cloud_region.dev.id
}

resource "cloud_stack_member" "default" {
  stack_id = cloud_stack.default.id
  user_id  = cloud_organization_member.default.user_id

  role = "ADMIN"
}
