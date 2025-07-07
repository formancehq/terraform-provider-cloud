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
  email = "test@formance.cloud"
  role  = "ADMIN"
}
