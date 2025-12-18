package resources_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/formancehq/formance-sdk-cloud-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-cloud-go/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"go.uber.org/mock/gomock"
)

func TestStackConfigure(t *testing.T) {
	test(t, func(ctx context.Context) {

		type testCase struct {
			providerData func(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) any
			expectedErr  error
		}

		for _, tc := range []testCase{
			{
				providerData: func(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) any {
					return any(nil)
				},
			},
			{
				providerData: func(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) any {
					return "something"
				},
				expectedErr: resources.ErrProviderDataNotSet,
			},
			{
				providerData: func(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) any {
					return internal.NewStore(sdkClient, tp)
				},
			},
		} {

			og := resources.NewStack()().(resource.ResourceWithConfigure)

			res := resource.ConfigureResponse{
				Diagnostics: []diag.Diagnostic{},
			}
			ctrl := gomock.NewController(t)
			tp := pkg.NewMockTokenProviderImpl(ctrl)
			apiMock := pkg.NewMockCloudSDK(ctrl)

			if tc.expectedErr == nil {
				tp.EXPECT().IntrospectToken(gomock.Any()).Return(oidc.IntrospectionResponse{
					Claims: map[string]interface{}{
						"organization_id": "organization_" + uuid.NewString()[:8],
					},
				}, nil).AnyTimes()
			}

			og.Configure(ctx, resource.ConfigureRequest{
				ProviderData: tc.providerData(apiMock, tp),
			}, &res)

			if tc.expectedErr != nil {
				require.Len(t, res.Diagnostics, 1, "Expected one diagnostic")
				require.Equal(t, res.Diagnostics[0].Summary(), tc.expectedErr.Error())
			} else {
				require.Empty(t, res.Diagnostics, "Expected no diagnostics")
			}

		}
	})
}

func TestStackMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		og := resources.NewStack()().(resource.ResourceWithConfigure)

		res := resource.MetadataResponse{}

		og.Metadata(ctx, resource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_stack")
	})
}

func TestStackCreate(t *testing.T) {
	type testCase struct {
		name     string
		regionID string
		version  string
	}

	for _, tc := range []testCase{
		{
			name: uuid.NewString(),
		},
		{},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			test(t, func(ctx context.Context) {
				r := resources.NewStack()().(resource.ResourceWithConfigure)
				organizationId := uuid.NewString()
				configureRes := resource.ConfigureResponse{
					Diagnostics: []diag.Diagnostic{},
				}
				ctrl := gomock.NewController(t)
				tp := pkg.NewMockTokenProviderImpl(ctrl)
				apiMock := pkg.NewMockCloudSDK(ctrl)
				store := internal.NewStore(apiMock, tp)

				tp.EXPECT().IntrospectToken(gomock.Any()).Return(oidc.IntrospectionResponse{
					Claims: map[string]interface{}{
						"organization_id": organizationId,
					},
				}, nil).AnyTimes()

				r.Configure(ctx, resource.ConfigureRequest{
					ProviderData: store,
				}, &configureRes)

				require.Empty(t, configureRes.Diagnostics, "Expected no diagnostics on configure")

				md := map[string]string{
					"github.com/formancehq/terraform-provider-cloud/protected": "true",
				}
				stackID := uuid.NewString()
				now := time.Now()
				apiMock.EXPECT().CreateStack(gomock.Any(), organizationId, &shared.CreateStackRequest{
					Name:     tc.name,
					Metadata: md,
					RegionID: tc.regionID,
					Version:  pointer.For(tc.version),
				}).Return(&operations.CreateStackResponse{
					StatusCode:  http.StatusCreated,
					RawResponse: &http.Response{StatusCode: http.StatusCreated},
					CreateStackResponse: &shared.CreateStackResponse{
						Data: &shared.Stack{
							ID:                       stackID,
							Name:                     tc.name,
							OrganizationID:           organizationId,
							RegionID:                 tc.regionID,
							Version:                  pointer.For(tc.version),
							URI:                      "https://example.com",
							Metadata:                 md,
							Status:                   shared.StackStatusReady,
							State:                    shared.StackStateActive,
							ExpectedStatus:           shared.ExpectedStatusReady,
							LastStateUpdate:          now,
							LastExpectedStatusUpdate: now,
							LastStatusUpdate:         now,
							Reachable:                true,
							StargateEnabled:          false,
							Synchronised:             true,
							Modules:                  []shared.Module{},
						},
					},
				}, nil)

				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: getSchemaTypes(resources.SchemaStack),
						}, map[string]tftypes.Value{
							"id":            tftypes.NewValue(tftypes.String, nil),
							"name":          tftypes.NewValue(tftypes.String, tc.name),
							"region_id":     tftypes.NewValue(tftypes.String, tc.regionID),
							"version":       tftypes.NewValue(tftypes.String, tc.version),
							"force_destroy": tftypes.NewValue(tftypes.Bool, nil),
							"uri":           tftypes.NewValue(tftypes.String, "https://example.com"),
							"metadata": tftypes.NewValue(tftypes.Map{
								ElementType: tftypes.String,
							}, nil),
						}),
						Schema: resources.SchemaStack,
					},
				}
				res := resource.CreateResponse{
					Diagnostics: []diag.Diagnostic{},
					State: tfsdk.State{
						Schema: resources.SchemaStack,
					},
				}
				r.Create(ctx, req, &res)

				require.Len(t, res.Diagnostics, 0, "Expected no diagnostics on create")

				model := &resources.StackModel{}
				res.State.Get(ctx, model)

			})
		})
	}
}

func TestStackValidateConfig(t *testing.T) {
	type testCase struct {
		organizationID *string
		regionID       *string
	}

	for _, tc := range []testCase{
		{},
		{
			organizationID: pointer.For(uuid.NewString()),
			regionID:       pointer.For(uuid.NewString()),
		},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			test(t, func(ctx context.Context) {
				r := resources.NewStack()().(resource.ResourceWithValidateConfig)

				res := resource.ValidateConfigResponse{
					Diagnostics: []diag.Diagnostic{},
				}

				r.ValidateConfig(ctx, resource.ValidateConfigRequest{
					Config: tfsdk.Config{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: getSchemaTypes(resources.SchemaStack),
						}, map[string]tftypes.Value{
							"name":          tftypes.NewValue(tftypes.String, nil),
							"region_id":     tftypes.NewValue(tftypes.String, tc.regionID),
							"version":       tftypes.NewValue(tftypes.String, nil),
							"id":            tftypes.NewValue(tftypes.String, nil),
							"force_destroy": tftypes.NewValue(tftypes.Bool, nil),
							"uri":           tftypes.NewValue(tftypes.String, nil),
							"metadata": tftypes.NewValue(tftypes.Map{
								ElementType: tftypes.String,
							}, nil),
						}),
						Schema: resources.SchemaStack,
					},
				}, &res)

				if tc.organizationID == nil || tc.regionID == nil {
					require.Len(t, res.Diagnostics, 1, "Expected one diagnostic for each missing field")
					require.Equal(t, res.Diagnostics[0].Summary(), "Invalid Region ID")

				} else {
					require.Empty(t, res.Diagnostics, "Expected no diagnostics")
				}
			})
		})
	}
}
