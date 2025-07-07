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

# TF_VAR_region_datasource_name
variable "region_datasource_name" {
  type = string
}


data "cloud_regions" "dev" {
  name = var.region_datasource_name
}

output "region_datasource_id" {
  value = data.cloud_regions.dev.id
}
