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
  type    = string
}

output "region_datasource_name" {
  value = data.formancecloud_regions.dev.id
}

data "formancecloud_regions" "dev" {
  name = var.region_datasource_name
  organization_id = formancecloud_organization.default.id
}

data "formancecloud_region_versions" "dev" {
  id = data.formancecloud_regions.dev.id
  organization_id = formancecloud_organization.default.id
}