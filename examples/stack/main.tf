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

# TF_VAR_region_datasource_id
variable "region_datasource_id" {
  type = string
}

import {
  to = cloud_region.dev
  id = var.region_datasource_id
}

resource "cloud_region" "dev" {}


resource "cloud_stack" "default" {
  name      = "test-stack"
  region_id = cloud_region.dev.id
}
