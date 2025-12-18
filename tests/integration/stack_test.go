package integration_test

import (
	"errors"
	"fmt"
	"net/http"
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

func TestStack(t *testing.T) {
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
					resource "cloud_stack" "test" {
						name = "test"
						region_id = "staging"
						metadata = {
							"env" = "test"
						}
						force_destroy = true
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
				stackID := uuid.NewString()
				md := map[string]string{
					"env": "test",
					"github.com/formancehq/terraform-provider-cloud/protected": "true",
				}
				now := time.Now()
				stackData := &shared.Stack{
					ID:             stackID,
					Name:           "test",
					OrganizationID: organizationID,
					RegionID:       "staging",
					Version:        pointer.For("latest"),
					URI:            "https://example.com",
					Metadata:       md,
					Status:         shared.StackStatusReady,
					State:          shared.StackStateActive,
					ExpectedStatus: shared.ExpectedStatusReady,
					LastStateUpdate:          now,
					LastExpectedStatusUpdate: now,
					LastStatusUpdate:         now,
					Reachable:                true,
					StargateEnabled:          false,
					Synchronised:             true,
					Modules:                  []shared.Module{},
				}
				cloudSdk.EXPECT().CreateStack(gomock.Any(), organizationID, gomock.Any()).
					Return(&operations.CreateStackResponse{
						StatusCode:  http.StatusCreated,
						RawResponse: &http.Response{StatusCode: http.StatusCreated},
						CreateStackResponse: &shared.CreateStackResponse{
							Data: stackData,
						},
					}, nil)
				cloudSdk.EXPECT().ReadStack(gomock.Any(), organizationID, stackID).
					Return(&operations.GetStackResponse{
						StatusCode:  http.StatusOK,
						RawResponse: &http.Response{StatusCode: http.StatusOK},
						CreateStackResponse: &shared.CreateStackResponse{
							Data: stackData,
						},
					}, nil)
				cloudSdk.EXPECT().DeleteStack(gomock.Any(), organizationID, stackID, true).Return(&operations.DeleteStackResponse{
					StatusCode:  http.StatusNoContent,
					RawResponse: &http.Response{StatusCode: http.StatusNoContent},
				}, nil)
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

func TestStackAlreadyDeleted(t *testing.T) {
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
					resource "cloud_stack" "test" {
						name = "test"
						region_id = "staging"
						metadata = {
							"env" = "test"
						}
						force_destroy = true
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
				stackID := uuid.NewString()
				md := map[string]string{
					"env": "test",
					"github.com/formancehq/terraform-provider-cloud/protected": "true",
				}
				now := time.Now()
				stackData := &shared.Stack{
					ID:             stackID,
					Name:           "test",
					OrganizationID: organizationID,
					RegionID:       "staging",
					Version:        pointer.For("latest"),
					URI:            "https://example.com",
					Metadata:       md,
					Status:         shared.StackStatusReady,
					State:          shared.StackStateActive,
					ExpectedStatus: shared.ExpectedStatusReady,
					LastStateUpdate:          now,
					LastExpectedStatusUpdate: now,
					LastStatusUpdate:         now,
					Reachable:                true,
					StargateEnabled:          false,
					Synchronised:             true,
					Modules:                  []shared.Module{},
				}
				cloudSdk.EXPECT().CreateStack(gomock.Any(), organizationID, gomock.Any()).
					Return(&operations.CreateStackResponse{
						StatusCode:  http.StatusCreated,
						RawResponse: &http.Response{StatusCode: http.StatusCreated},
						CreateStackResponse: &shared.CreateStackResponse{
							Data: stackData,
						},
					}, nil)
				cloudSdk.EXPECT().ReadStack(gomock.Any(), organizationID, stackID).
					Return(&operations.GetStackResponse{
						StatusCode:  http.StatusOK,
						RawResponse: &http.Response{StatusCode: http.StatusOK},
						CreateStackResponse: &shared.CreateStackResponse{
							Data: stackData,
						},
					}, nil)
				cloudSdk.EXPECT().DeleteStack(gomock.Any(), organizationID, stackID, true).Return(&operations.DeleteStackResponse{
					StatusCode:  http.StatusNotFound,
					RawResponse: &http.Response{StatusCode: http.StatusNotFound},
				}, errors.New("stack not found"))
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
