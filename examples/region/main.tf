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

resource "formancecloud_region" "dev" {
  depends_on = [formancecloud_organization.default]
  name       = "dev"
}
