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

func TestStackModule(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"cloud": providerserver.NewProtocol6WithError(Provider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "cloud" {}
					data "cloud_regions" "dev" {
						name = "%s"
					}

					output "region_id" {
						value = data.cloud_regions.dev.id
					}

					resource "cloud_stack" "default" {
						name = "test"
						region_id = data.cloud_regions.dev.id

						version = "default"
						force_destroy = true
					}

					resource "cloud_stack_module" "default_webhooks" {
						name = "webhooks"
						stack_id = cloud_stack.default.id
					}

					resource "cloud_stack_module" "default_reconciliation" {
						name = "reconciliation"
						stack_id = cloud_stack.default.id
					}

					resource "cloud_stack_module" "default_orchestration" {
						name = "orchestration"
						stack_id = cloud_stack.default.id
					}

				`, RegionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchOutput("region_id", regexp.MustCompile(`.+`)),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("cloud_stack.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("cloud_stack.default", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
					statecheck.ExpectKnownValue("cloud_stack.default", tfjsonpath.New("force_destroy"), knownvalue.Bool(true)),
				},
			},
			{
				RefreshState: true,
			},
		},
	})

}
