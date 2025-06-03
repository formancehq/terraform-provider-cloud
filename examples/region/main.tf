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

resource "formancecloud_region" "dev" {
  depends_on      = [formancecloud_organization.default]
  name            = "dev"
  organization_id = formancecloud_organization.default.id
}
