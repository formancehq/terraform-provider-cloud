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
  name = var.region_datasource_name
}

data "formancecloud_region_versions" "dev" {
  id = data.formancecloud_regions.dev.id
}