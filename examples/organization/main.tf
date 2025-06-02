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
