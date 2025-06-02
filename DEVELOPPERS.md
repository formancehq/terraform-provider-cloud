# Formance Terraform Provider

Create a `~/.terraformrc` with the following content:

replace {WORKING_DIRECTORY} with your username

```hcl
provider_installation {
  dev_overrides {
    "formancehq/tf-cloud-provider" = "${WORKING_DIRECTORY}/build"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
