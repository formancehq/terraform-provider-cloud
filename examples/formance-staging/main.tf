terraform {
  required_providers {
    cloud = {
      source = "formancehq/cloud"
    }
  }
}

provider "cloud" {
  endpoint = "https://app.staging.formance.cloud/api"
}

resource "cloud_organization" "default" {
  name = "formancehq"
}
