terraform {
  required_providers {
    formancecloud = {
      source = "formancehq/cloud"
    }
  }
}

provider "formancecloud" {
  endpoint = "https://app.staging.formance.cloud/api"
}

resource "formancecloud_organization" "default" {
  name = "formancehq"
}
