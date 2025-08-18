// Format Terraform code for use in documentation.
// If you do not have Terraform installed, you can remove the formatting command, but it is suggested
// to ensure the documentation is formatted properly.

// //go:generate rm -rf docs
// //go:generate mkdir docs
// //go:generate terraform fmt -recursive ./examples/

// // Generate documentation.
// //go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest generate --provider-dir . -provider-name cloud

// Generate terraform providers schema
//
//go:generate rm -rf ./schema && mkdir -p ./schema && cd examples/latest && terraform init -upgrade && terraform providers schema -json > ../../schema/providers.json
package main

import "github.com/formancehq/terraform-provider-cloud/cmd"

func main() {
	cmd.Execute()
}
