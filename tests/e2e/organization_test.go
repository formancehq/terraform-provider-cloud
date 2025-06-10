package e2e_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var OrganizationId = ""

func TestOrganization(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"formancecloud": providerserver.NewProtocol6WithError(Provider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "formancecloud" {}

					import {
						id = "` + OrganizationId + `"
						to = formancecloud_organization.default
					}

					resource "formancecloud_organization" "default" {
						name = "default"
					}

				`,
			},
			{
				Config: `
					provider "formancecloud" {}

					import {
						id = "` + OrganizationId + `"
						to = formancecloud_organization.default
					}

					resource "formancecloud_organization" "default" {
						name = "newName"
						domain = "newDomain"
						default_organization_access = "ADMIN"
						default_stack_access = "ADMIN"
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("formancecloud_organization.default", tfjsonpath.New("name"), knownvalue.StringExact("newName")),
					statecheck.ExpectKnownValue("formancecloud_organization.default", tfjsonpath.New("domain"), knownvalue.StringExact("newDomain")),
					statecheck.ExpectKnownValue("formancecloud_organization.default", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("formancecloud_organization.default", tfjsonpath.New("default_organization_access"), knownvalue.StringExact("ADMIN")),
					statecheck.ExpectKnownValue("formancecloud_organization.default", tfjsonpath.New("default_stack_access"), knownvalue.StringExact("ADMIN")),
				},
			},
			{
				RefreshState: true,
			},
		},
	})

}
