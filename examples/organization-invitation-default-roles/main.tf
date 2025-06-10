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
  name                        = "formancehq"
  domain                      = "exemple.com"
  default_organization_access = "ADMIN"
  default_stack_access        = "ADMIN"
}
