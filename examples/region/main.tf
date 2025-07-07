terraform {
  required_providers {
    cloud = {
      source = "formancehq/cloud"
      configuration_aliases = [
        cloud.dev,
        cloud.prod
      ]
      version = "0.0.5"
    }
  }
}

provider "cloud" {
  alias = "dev"
}

provider "cloud" {
  alias = "prod"
}

data "cloud_current_organization" "default" {
  provider = cloud.prod
}

data "cloud_regions" "dev" {
  provider = cloud.dev
  name     = "staging"
}
