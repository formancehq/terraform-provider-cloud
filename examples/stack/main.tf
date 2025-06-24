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

# TF_VAR_region_datasource_id
variable "region_datasource_id" {
  type = string
}

import {
  to = formancecloud_region.dev
  id = var.region_datasource_id
}

resource "formancecloud_region" "dev" {}


resource "formancecloud_stack" "default" {
  name      = "test-stack"
  region_id = formancecloud_region.dev.id
}
