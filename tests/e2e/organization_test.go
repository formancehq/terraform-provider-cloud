package e2e_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestOrganization(t *testing.T) {
	t.Parallel()
	domain := uuid.NewString()
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
					resource "formancecloud_organization" "test" {
						name = "formancehq"
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("name"), knownvalue.StringExact("formancehq")),
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("domain"), knownvalue.Null()),
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("default_organization_access"), knownvalue.StringExact("GUEST")),
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("default_stack_access"), knownvalue.StringExact("GUEST")),
				},
			},
			{
				Config: fmt.Sprintf(`
					provider "formancecloud" {}
					resource "formancecloud_organization" "test" {
						name = "newName"
						domain = "%s"
						default_organization_access = "ADMIN"
						default_stack_access = "ADMIN"
					}
				`, domain),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("name"), knownvalue.StringExact("newName")),
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("domain"), knownvalue.StringExact(domain)),
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("default_organization_access"), knownvalue.StringExact("ADMIN")),
					statecheck.ExpectKnownValue("formancecloud_organization.test", tfjsonpath.New("default_stack_access"), knownvalue.StringExact("ADMIN")),
				},
			},
			{
				RefreshState: true,
			},
			{
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				ResourceName:    "formancecloud_organization.test",
			},
		},
	})

}
