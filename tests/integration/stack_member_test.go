package integration_test

import (
	"fmt"
	"net/http"
	"regexp"
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
						provider "formancecloud" {}
						resource "formancecloud_stack_member" "test" {
							user_id  = "user-id-123"
							stack_id = "stack-id-456"
							role	 = "GUEST"
						}
					`,
				},
				{
					Config: `
						provider "formancecloud" {}
						resource "formancecloud_stack_member" "test" {
							user_id  = "user-id-123"
							stack_id = "stack-id-456"
							role	 = "ADMIN"
						}
					`,
				},
			},
			expectedCalls: func(cloudSdk *pkg.MockCloudSDK, tokenProvider *pkg.MockTokenProviderImpl) {
				cloudSdk.EXPECT().UpsertStackUserAccess(gomock.Any(), "client_id", "stack-id-456", "user-id-123", sdk.UpdateStackUserRequest{
					Role: sdk.GUEST,
				}).Return(nil, nil)
				cloudSdk.EXPECT().ReadStackUserAccess(gomock.Any(), "client_id", "stack-id-456", "user-id-123").
					Return(&sdk.ReadStackUserAccess{
						Data: &sdk.StackUserAccess{
							StackId: "stack-id-456",
							UserId:  "user-id-123",
							Email:   "example@formance.com",
							Role:    sdk.GUEST,
						},
					}, nil, nil).Times(2)
				cloudSdk.EXPECT().UpsertStackUserAccess(gomock.Any(), "client_id", "stack-id-456", "user-id-123", sdk.UpdateStackUserRequest{
					Role: sdk.ADMIN,
				}).Return(nil, nil)
				cloudSdk.EXPECT().ReadStackUserAccess(gomock.Any(), "client_id", "stack-id-456", "user-id-123").
					Return(&sdk.ReadStackUserAccess{
						Data: &sdk.StackUserAccess{
							StackId: "stack-id-456",
							UserId:  "user-id-123",
							Email:   "example@formance.com",
							Role:    sdk.ADMIN,
						},
					}, nil, nil)
				cloudSdk.EXPECT().DeleteStackUserAccess(gomock.Any(), "client_id", "stack-id-456", "user-id-123").Return(nil, nil)
			},
		},
		{
			step: []resource.TestStep{
				{
					Config: `
						provider "formancecloud" {}
						resource "formancecloud_stack_member" "test" {
							user_id  = "user-id-123"
							stack_id = "stack-id-456"
						}
					`,
					ExpectError: regexp.MustCompile(`The argument "role" is required`),
				},
			},
		},
	} {

		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cloudSdk := pkg.NewMockCloudSDK(ctrl)
			tokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
			cloudProvider := server.NewProvider(
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

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"formancecloud": providerserver.NewProtocol6WithError(cloudProvider()),
				},
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version0_15_0),
				},
				Steps: tc.step,
			})
		})
	}
}
