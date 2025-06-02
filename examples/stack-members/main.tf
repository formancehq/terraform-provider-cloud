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
  email = "example@formance.dev"
  role = "ADMIN"
}

resource "formancecloud_region" "dev" {}

resource "formancecloud_stack" "default" {
  name = "default"
  organization_id = formancecloud_organization.default.id
  region_id = formancecloud_region.dev.id
}

resource "formancecloud_stack_member" "default" {
  stack_id = formancecloud_stack.default.id
  organization_id = formancecloud_organization.default.id
  user_id = formancecloud_organization_member.default.user_id
  
  role = "ADMIN"
}
