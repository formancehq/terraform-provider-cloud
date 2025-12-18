package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/operations"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"go.opentelemetry.io/otel/trace/noop"
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
			expectedCalls: func(cloudSdk *pkg.MockCloudSDK, tokenProvider *pkg.MockTokenProviderImpl) {
				organizationID := uuid.NewString()
				tokenProvider.EXPECT().IntrospectToken(gomock.Any()).Return(oidc.IntrospectionResponse{
					Claims: map[string]interface{}{
						"organization_id": organizationID,
					},
				}, nil).AnyTimes()
				cloudSdk.EXPECT().GetRegionVersions(gomock.All(), organizationID, "another-region-id").Return(&operations.GetRegionVersionsResponse{
					StatusCode:  http.StatusOK,
					RawResponse: &http.Response{StatusCode: http.StatusOK},
					GetRegionVersionsResponse: &shared.GetRegionVersionsResponse{
						Data: []shared.Version{
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
					},
				}, nil).AnyTimes()
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
				organizationID := uuid.NewString()
				tokenProvider.EXPECT().IntrospectToken(gomock.Any()).Return(oidc.IntrospectionResponse{
					Claims: map[string]interface{}{
						"organization_id": organizationID,
					},
				}, nil).AnyTimes()
				cloudSdk.EXPECT().ListRegions(gomock.Any(), organizationID).Return(&operations.ListRegionsResponse{
					StatusCode:  http.StatusOK,
					RawResponse: &http.Response{StatusCode: http.StatusOK},
					ListRegionsResponse: &shared.ListRegionsResponse{
						Data: []shared.AnyRegion{
							{
								ID:   "some-region-id",
								Name: "Some Region",
							},
							{ // <- This region will be sorted in first position
								ID:   "another-region-id",
								Name: "Another Region",
							},
						},
					},
				}, nil).AnyTimes()

				cloudSdk.EXPECT().GetRegionVersions(gomock.All(), organizationID, "another-region-id").Return(&operations.GetRegionVersionsResponse{
					StatusCode:  http.StatusOK,
					RawResponse: &http.Response{StatusCode: http.StatusOK},
					GetRegionVersionsResponse: &shared.GetRegionVersionsResponse{
						Data: []shared.Version{
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
					},
				}, nil).AnyTimes()
			},
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			cloudSdk := pkg.NewMockCloudSDK(ctrl)
			tokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
			cloudProvider := server.NewProvider(
				noop.NewTracerProvider(),

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

			resource.ParallelTest(t, resource.TestCase{
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
