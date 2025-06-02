terraform {
  required_providers {
    formancecloud = {
      source = "formancehq/tf-cloud-provider"
    }
  }
}

provider "formancecloud" {}

resource "formancecloud_organization" "default" {
    name = "default"
}

resource "formancecloud_organization_member" "default" {
    organization_id = formancecloud_organization.default.id
    email = "test@formance.cloud"
    role = "ADMIN"
}
