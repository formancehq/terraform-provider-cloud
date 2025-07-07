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

variable "region_datasource_name" {
  type = string
}

output "region_datasource_name" {
  value = data.cloud_regions.dev.id
}

data "cloud_regions" "dev" {
  name = var.region_datasource_name
}

data "cloud_region_versions" "dev" {
  id = data.cloud_regions.dev.id
}