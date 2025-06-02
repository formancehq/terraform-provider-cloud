terraform {
  required_providers {
    formancecloud = {
      source = "formancehq/tf-cloud-provider"
    }
  }
}

provider "formancecloud" {}

resource "formancecloud_organization" "default" {
  name = "formancehq"
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
  name = "test-stack"
  organization_id = formancecloud_organization.default.id
  region_id = formancecloud_region.dev.id
}
