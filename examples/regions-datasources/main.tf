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

# TF_VAR_region_datasource_name
variable "region_datasource_name" {
  type = string
}


data "formancecloud_regions" "dev" {
  name            = var.region_datasource_name
  organization_id = formancecloud_organization.default.id
}

output "region_datasource_id" {
  value = data.formancecloud_regions.dev.id
}
