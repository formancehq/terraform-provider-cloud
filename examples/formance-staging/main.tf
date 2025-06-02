terraform {
  required_providers {
    formancecloud = {
      source = "formancehq/cloud"
    }
  }
}

provider "formancecloud" {
  endpoint     = "https://app.staging.formance.cloud/api"
}

var "staging_region_id" {
  type = string
}

import {
  to = formancecloud_region.region.public
  id = var.staging_region_id
}

resource "formancecloud_region" "region" "public" {}


resource "formancecloud_organization" "default" {
  name = "formancehq"
}

resource "formancecloud_region" "region" "private" {
    name       = "eu-west-data-plane"
    depends_on = [formancecloud_organization.default]
    organization_id = formancecloud_organization.default.id
}

resource "formancecloud_region" "region" "private" {
    name       = "us-east-data-plane"
    depends_on = [formancecloud_organization.default]
    organization_id = formancecloud_organization.default.id
}
