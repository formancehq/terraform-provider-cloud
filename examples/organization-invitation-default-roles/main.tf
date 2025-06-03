terraform {
  required_providers {
    formancecloud = {
      source = "formancehq/cloud"
    }
  }
}

provider "formancecloud" {}

resource "formancecloud_organization" "default" {
  name                        = "formancehq"
  domain                      = "exemple.com"
  default_organization_access = "ADMIN"
  default_stack_access        = "ADMIN"
}
