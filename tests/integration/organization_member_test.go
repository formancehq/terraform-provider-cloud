package integration_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/formancehq/formance-sdk-cloud-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-cloud-go/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"go.opentelemetry.io/otel/trace/noop"
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
				invitationID := uuid.NewString()
				now := time.Now()
				invitation := &shared.Invitation{
					ID:             invitationID,
					UserEmail:      "example@formance.com",
					Status:         shared.InvitationStatusPending,
					OrganizationID: organizationID,
					CreationDate:   now,
				}

				mcs.EXPECT().
					CreateInvitation(gomock.Any(), organizationID, "example@formance.com").
					Return(&operations.CreateInvitationResponse{
						StatusCode:  http.StatusCreated,
						RawResponse: &http.Response{StatusCode: http.StatusCreated},
						CreateInvitationResponse: &shared.CreateInvitationResponse{
							Data: invitation,
						},
					}, nil)

				mcs.EXPECT().ListOrganizationInvitations(gomock.Any(), organizationID).Return(
					&operations.ListInvitationsResponse{
						StatusCode:  http.StatusOK,
						RawResponse: &http.Response{StatusCode: http.StatusOK},
						ListInvitationsResponse: &shared.ListInvitationsResponse{
							Data: []shared.Invitation{
								*invitation,
							},
						},
					},
					nil,
				).AnyTimes()

				mcs.EXPECT().DeleteInvitation(gomock.Any(), organizationID, invitationID).Return(
					&operations.DeleteInvitationResponse{
						StatusCode:  http.StatusNoContent,
						RawResponse: &http.Response{StatusCode: http.StatusNoContent},
					}, nil,
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
				invitationID := uuid.NewString()
				userID := uuid.NewString()
				now := time.Now()
				invitation := &shared.Invitation{
					ID:             invitationID,
					UserEmail:      "example@formance.com",
					Status:         shared.InvitationStatusAccepted,
					OrganizationID: organizationID,
					UserID:         pointer.For(userID),
					CreationDate:   now,
				}

				mcs.EXPECT().
					CreateInvitation(gomock.Any(), organizationID, "example@formance.com").
					Return(&operations.CreateInvitationResponse{
						StatusCode:  http.StatusCreated,
						RawResponse: &http.Response{StatusCode: http.StatusCreated},
						CreateInvitationResponse: &shared.CreateInvitationResponse{
							Data: invitation,
						},
					}, nil)

				mcs.EXPECT().ListOrganizationInvitations(gomock.Any(), organizationID).Return(
					&operations.ListInvitationsResponse{
						StatusCode:  http.StatusOK,
						RawResponse: &http.Response{StatusCode: http.StatusOK},
						ListInvitationsResponse: &shared.ListInvitationsResponse{
							Data: []shared.Invitation{
								*invitation,
							},
						},
					},
					nil,
				).AnyTimes()

				mcs.EXPECT().ReadUserOfOrganization(gomock.Any(), organizationID, userID).Return(
					&operations.ReadUserOfOrganizationResponse{
						StatusCode:  http.StatusOK,
						RawResponse: &http.Response{StatusCode: http.StatusOK},
						ReadOrganizationUserResponse: &shared.ReadOrganizationUserResponse{
							Data: &shared.ReadOrganizationUserResponseData{
								Email:    "example@formance.com",
								ID:       userID,
								PolicyID: 0,
							},
						},
					}, nil,
				)

				mcs.EXPECT().DeleteUserOfOrganization(gomock.Any(), organizationID, userID).Return(
					&operations.DeleteUserFromOrganizationResponse{
						StatusCode:  http.StatusNoContent,
						RawResponse: &http.Response{StatusCode: http.StatusNoContent},
					}, nil,
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
