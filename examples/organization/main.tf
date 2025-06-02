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
