package integration_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"go.uber.org/mock/gomock"
)

func TestOrganizationMember(t *testing.T) {
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

						resource "cloud_organization_member" "default" {
							email = "example@formance.com"
						}
					`,
				},
			},
			expectedCalls: func(mcs *pkg.MockCloudSDK, mtpi *pkg.MockTokenProviderImpl) {
				organizationID := uuid.NewString()
				res := oidc.IntrospectionResponse{
					Claims: map[string]interface{}{
						"organization_id": organizationID,
					},
				}
				mtpi.EXPECT().IntrospectToken(gomock.Any()).Return(res, nil).AnyTimes()
				invitation := &sdk.Invitation{
					Id:        uuid.NewString(),
					Role:      "",
					UserEmail: "example@formance.com",
				}

				mcs.EXPECT().
					CreateInvitation(gomock.Any(), organizationID, "example@formance.com", sdk.InvitationClaim{}).
					Return(&sdk.CreateInvitationResponse{
						Data: invitation,
					}, nil, nil)

				invitation.Status = "PENDING"
				mcs.EXPECT().ListOrganizationInvitations(gomock.Any(), organizationID).Return(
					&sdk.ListInvitationsResponse{
						Data: []sdk.Invitation{
							*invitation,
						},
					},
					nil, nil,
				).AnyTimes()

				mcs.EXPECT().DeleteInvitation(gomock.Any(), organizationID, invitation.Id).Return(
					nil, nil,
				)
			},
		},
		{
			step: []resource.TestStep{
				{
					Config: `
						provider "cloud" {}

						resource "cloud_organization_member" "default" {
							email = "example@formance.com"
						}
					`,
				},
			},
			expectedCalls: func(mcs *pkg.MockCloudSDK, mtpi *pkg.MockTokenProviderImpl) {
				organizationID := uuid.NewString()
				res := oidc.IntrospectionResponse{
					Claims: map[string]interface{}{
						"organization_id": organizationID,
					},
				}
				mtpi.EXPECT().IntrospectToken(gomock.Any()).Return(res, nil).AnyTimes()
				invitation := &sdk.Invitation{
					Id:        uuid.NewString(),
					Role:      "",
					UserEmail: "example@formance.com",
				}

				mcs.EXPECT().
					CreateInvitation(gomock.Any(), organizationID, "example@formance.com", sdk.InvitationClaim{}).
					Return(&sdk.CreateInvitationResponse{
						Data: invitation,
					}, nil, nil)

				invitation.Status = "ACCEPTED"
				invitation.UserId = pointer.For(uuid.NewString())
				mcs.EXPECT().ListOrganizationInvitations(gomock.Any(), organizationID).Return(
					&sdk.ListInvitationsResponse{
						Data: []sdk.Invitation{
							*invitation,
						},
					},
					nil, nil,
				).AnyTimes()

				mcs.EXPECT().ReadUserOfOrganization(gomock.Any(), organizationID, *invitation.UserId).Return(
					&sdk.ReadOrganizationUserResponse{
						Data: &sdk.OrganizationUser{
							Role:  "GUEST",
							Email: "example@formance.com",
							Id:    *invitation.UserId,
						},
					}, nil, nil,
				)

				mcs.EXPECT().DeleteUserOfOrganization(gomock.Any(), organizationID, *invitation.UserId).Return(
					nil, nil,
				)
			},
		},
		{
			step: []resource.TestStep{
				{
					Config: `
						provider "cloud" {}

						resource "cloud_organization_member" "default" {}
					`,
					ExpectError: regexp.MustCompile(`"email" is required`),
				},
			},
		},
	} {
		t.Run(t.Name(), func(t *testing.T) {
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
