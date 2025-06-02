package e2e_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// On dev env search for name = "https://${USER}.formance.dev"
// On staging env search for name = "staging"
// On production env search for name = "production"

var (
	RegionName = ""
)

func TestRegionDatasources(t *testing.T) {
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
					resource "formancecloud_organization" "default" {
						name = "formancehq"
					}

					data "formancecloud_regions" "dev" {
						name = "%s"
						organization_id = formancecloud_organization.default.id
					}

					output "region_name" {
						value = data.formancecloud_regions.dev.name
					}
					output "region_id" {
						value = data.formancecloud_regions.dev.id
					}
					output "region_organization_id" {
						value = data.formancecloud_regions.dev.organization_id
					}
				`, RegionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("region_name", RegionName),
					resource.TestMatchOutput("region_id", regexp.MustCompile(`.+`)),
				),
			},
			{
				RefreshState: true,
			},
		},
	})
}
