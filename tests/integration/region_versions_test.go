package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"go.uber.org/mock/gomock"
)

func TestRegionVersions(t *testing.T) {
	t.Parallel()

	type testCase struct {
		step          resource.TestStep
		expectedCalls func(*pkg.MockCloudSDK, *pkg.MockTokenProviderImpl)
	}

	for i, tc := range []testCase{
		{
			step: resource.TestStep{
				Config: `
					provider "cloud" {}

					data "cloud_region_versions" "default" {
						id = "another-region-id"
					}
				`,
			},
		},
		{
			step: resource.TestStep{
				Config: `
				provider "cloud" {}
				
				data "cloud_region_versions" "default" {}				
				`,
			},
			expectedCalls: func(cloudSdk *pkg.MockCloudSDK, tokenProvider *pkg.MockTokenProviderImpl) {
				cloudSdk.EXPECT().ListRegions(gomock.Any(), "client_id").Return(&sdk.ListRegionsResponse{
					Data: []sdk.AnyRegion{
						{
							Id:   "some-region-id",
							Name: "Some Region",
						},
						{ // <- This region will be sorted in first position
							Id:   "another-region-id",
							Name: "Another Region",
						},
					},
				}, nil, nil).AnyTimes()
			},
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			cloudSdk := pkg.NewMockCloudSDK(ctrl)
			tokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
			cloudProvider := server.NewProvider(
				logging.Testing().WithField("test", fmt.Sprintf("region_versions_test_%d", i)),
				server.FormanceCloudEndpoint("dummy-endpoint"),
				server.FormanceCloudClientId("organization_client_id"),
				server.FormanceCloudClientSecret("dummy-client-secret"),
				transport,
				func(creds pkg.Creds, transport http.RoundTripper) pkg.CloudSDK {
					return cloudSdk
				},
				func(transport http.RoundTripper, creds pkg.Creds) pkg.TokenProviderImpl {
					return tokenProvider
				},
			)

			if tc.expectedCalls != nil {
				tc.expectedCalls(cloudSdk, tokenProvider)
			}
			cloudSdk.EXPECT().GetRegionVersions(gomock.All(), "client_id", "another-region-id").Return(&sdk.GetRegionVersionsResponse{
				Data: []sdk.Version{
					{
						Name: "v1.0.0",
						Versions: map[string]string{
							"ledger":   "v1.0.0",
							"payments": "v1.0.0",
						},
					},
					{
						Name: "v2.0.0",
						Versions: map[string]string{
							"ledger":   "v2.0.12",
							"payments": "v2.2.1",
						},
					},
				},
			}, nil, nil).AnyTimes()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"cloud": providerserver.NewProtocol6WithError(cloudProvider()),
				},
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version0_15_0),
				},
				Steps: []resource.TestStep{
					tc.step,
				},
			})
		})
	}

}
