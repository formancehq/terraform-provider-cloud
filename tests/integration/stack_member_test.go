package integration_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/operations"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
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

func TestStackMember(t *testing.T) {
	t.Parallel()
	type testCase struct {
		step          []resource.TestStep
		expectedCalls func(*pkg.MockCloudSDK, *pkg.MockTokenProviderImpl)
	}

	for i, tc := range []testCase{
		{
			step: []resource.TestStep{
				{
					Config: `
						provider "cloud" {}
						resource "cloud_stack_member" "test" {
							user_id  = "user-id-123"
							stack_id = "stack-id-456"
							policy_id = 1
						}
					`,
				},
				{
					Config: `
						provider "cloud" {}
						resource "cloud_stack_member" "test" {
							user_id  = "user-id-123"
							stack_id = "stack-id-456"
							policy_id = 2
						}
					`,
				},
			},
			expectedCalls: func(cloudSdk *pkg.MockCloudSDK, tokenProvider *pkg.MockTokenProviderImpl) {
				organizationID := uuid.NewString()
				tokenProvider.EXPECT().IntrospectToken(gomock.Any()).Return(oidc.IntrospectionResponse{
					Claims: map[string]interface{}{
						"organization_id": organizationID,
					},
				}, nil).AnyTimes()

				cloudSdk.EXPECT().UpsertStackUserAccess(gomock.Any(), organizationID, "stack-id-456", "user-id-123", &shared.UpdateStackUserRequest{
					PolicyID: 1,
				}).Return(&operations.UpsertStackUserAccessResponse{
					StatusCode:  http.StatusOK,
					RawResponse: &http.Response{StatusCode: http.StatusOK},
				}, nil)
				cloudSdk.EXPECT().ReadStackUserAccess(gomock.Any(), organizationID, "stack-id-456", "user-id-123").
					Return(&operations.ReadStackUserAccessResponse{
						StatusCode:  http.StatusOK,
						RawResponse: &http.Response{StatusCode: http.StatusOK},
						ReadStackUserAccess: &shared.ReadStackUserAccess{
							Data: &shared.ReadStackUserAccessData{
								StackID:  "stack-id-456",
								UserID:   "user-id-123",
								Email:    "example@formance.com",
								PolicyID: 1,
							},
						},
					}, nil).Times(2)
				cloudSdk.EXPECT().UpsertStackUserAccess(gomock.Any(), organizationID, "stack-id-456", "user-id-123", &shared.UpdateStackUserRequest{
					PolicyID: 2,
				}).Return(&operations.UpsertStackUserAccessResponse{
					StatusCode:  http.StatusOK,
					RawResponse: &http.Response{StatusCode: http.StatusOK},
				}, nil)
				cloudSdk.EXPECT().ReadStackUserAccess(gomock.Any(), organizationID, "stack-id-456", "user-id-123").
					Return(&operations.ReadStackUserAccessResponse{
						StatusCode:  http.StatusOK,
						RawResponse: &http.Response{StatusCode: http.StatusOK},
						ReadStackUserAccess: &shared.ReadStackUserAccess{
							Data: &shared.ReadStackUserAccessData{
								StackID:  "stack-id-456",
								UserID:   "user-id-123",
								Email:    "example@formance.com",
								PolicyID: 2,
							},
						},
					}, nil)
				cloudSdk.EXPECT().DeleteStackUserAccess(gomock.Any(), organizationID, "stack-id-456", "user-id-123").Return(&operations.DeleteStackUserAccessResponse{
					StatusCode:  http.StatusNoContent,
					RawResponse: &http.Response{StatusCode: http.StatusNoContent},
				}, nil)
			},
		},
		{
			step: []resource.TestStep{
				{
					Config: `
						provider "cloud" {}
						resource "cloud_stack_member" "test" {
							user_id  = "user-id-123"
							stack_id = "stack-id-456"
						}
					`,
					ExpectError: regexp.MustCompile(`The argument "policy_id" is required`),
				},
			},
		},
	} {

		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			cloudSdk := pkg.NewMockCloudSDK(ctrl)
			tokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
			cloudProvider := server.NewProvider(
				noop.NewTracerProvider(),

				logging.Testing().WithField("test", fmt.Sprintf("test_%d", i)),
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
				Steps: tc.step,
			})
		})
	}
}
