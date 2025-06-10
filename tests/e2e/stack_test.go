package e2e_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestStack(t *testing.T) {
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
				Config: fmt.Sprintf(`
					provider "formancecloud" {}
					import {
						id = "`+OrganizationId+`"
						to = formancecloud_organization.default
					}

					resource "formancecloud_organization" "default" {
						name = "default"
					}

					data "formancecloud_regions" "dev" {
						name = "%s"
						organization_id = formancecloud_organization.default.id
					}

					output "region_id" {
						value = data.formancecloud_regions.dev.id
					}

					resource "formancecloud_stack" "test" {
						name = "test"
						organization_id = formancecloud_organization.default.id
						region_id = data.formancecloud_regions.dev.id

						force_destroy = true
					}
				`, RegionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchOutput("region_id", regexp.MustCompile(`.+`)),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("formancecloud_stack.test", tfjsonpath.New("name"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("formancecloud_stack.test", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
					statecheck.ExpectKnownValue("formancecloud_stack.test", tfjsonpath.New("force_destroy"), knownvalue.Bool(true)),
				},
			},
			{
				RefreshState: true,
			},
		},
	})

}
