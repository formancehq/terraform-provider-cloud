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
