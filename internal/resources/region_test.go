package resources_test

import (
	"context"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegionConfigure(t *testing.T) {
	test(t, func(ctx context.Context) {

		type testCase struct {
			providerData any
			expectedErr  error
		}

		for _, tc := range []testCase{
			{},
			{
				providerData: "something",
				expectedErr:  resources.ErrProviderDataNotSet,
			},
			{
				providerData: pkg.NewMockDefaultAPI(gomock.NewController(t)),
			},
		} {

			og := resources.NewRegion(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

			res := resource.ConfigureResponse{
				Diagnostics: []diag.Diagnostic{},
			}

			og.Configure(ctx, resource.ConfigureRequest{
				ProviderData: tc.providerData,
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

func TestRegionMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		og := resources.NewRegion(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

		res := resource.MetadataResponse{}

		og.Metadata(ctx, resource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_region")
	})
}

func TestRegionValidateConfig(t *testing.T) {
	type testCase struct {
		organizationID *string
	}

	for _, tc := range []testCase{
		{
			organizationID: pointer.For(uuid.NewString()),
		},
		{},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			test(t, func(ctx context.Context) {
				r := resources.NewRegion(logging.FromContext(ctx))().(resource.ResourceWithValidateConfig)

				res := resource.ValidateConfigResponse{
					Diagnostics: []diag.Diagnostic{},
				}

				r.ValidateConfig(ctx, resource.ValidateConfigRequest{
					Config: tfsdk.Config{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name":            tftypes.String,
								"organization_id": tftypes.String,
								"base_url":        tftypes.String,
								"secret":          tftypes.String,
								"id":              tftypes.String,
							},
						}, map[string]tftypes.Value{
							"name":            tftypes.NewValue(tftypes.String, nil),
							"organization_id": tftypes.NewValue(tftypes.String, tc.organizationID),
							"base_url":        tftypes.NewValue(tftypes.String, ""),
							"secret":          tftypes.NewValue(tftypes.String, ""),
							"id":              tftypes.NewValue(tftypes.String, nil),
						}),
						Schema: resources.SchemaRegion,
					},
				}, &res)

				if tc.organizationID == nil {
					require.Len(t, res.Diagnostics, 1, "Expected one diagnostic on validate config")
				} else {
					require.Empty(t, res.Diagnostics, "Expected no diagnostics on validate config")
				}

			})
		})
	}
}

func TestRegionCreate(t *testing.T) {
	type testCase struct {
		name           string
		organizationID string
	}

	for _, tc := range []testCase{
		{
			name:           uuid.NewString(),
			organizationID: uuid.NewString(),
		},
		{},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			test(t, func(ctx context.Context) {
				region := resources.NewRegion(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

				configureRes := resource.ConfigureResponse{
					Diagnostics: []diag.Diagnostic{},
				}
				ctrl := gomock.NewController(t)
				apiMock := pkg.NewMockDefaultAPI(ctrl)
				region.Configure(ctx, resource.ConfigureRequest{
					ProviderData: apiMock,
				}, &configureRes)

				require.Empty(t, configureRes.Diagnostics, "Expected no diagnostics on configure")

				apiMock.EXPECT().CreatePrivateRegion(gomock.Any(), tc.organizationID).Return(sdk.ApiCreatePrivateRegionRequest{
					ApiService: apiMock,
				})

				expectedBody := sdk.CreatePrivateRegionRequest{
					Name: tc.name,
				}

				apiMock.EXPECT().CreatePrivateRegionExecute(
					gomock.Cond(
						func(r sdk.ApiCreatePrivateRegionRequest) bool {
							return cmp.Equal(r.GetCreatePrivateRegionRequest(), &expectedBody, cmp.AllowUnexported(sdk.CreatePrivateRegionRequest{}), cmpopts.IgnoreUnexported(sdk.CreatePrivateRegionRequest{}))
						}),
				).Return(&sdk.CreatedPrivateRegionResponse{
					Data: sdk.PrivateRegion{
						Id:             uuid.NewString(),
						OrganizationID: tc.organizationID,
						Name:           tc.name,
						BaseUrl:        "https://example.com",
						Secret: &sdk.PrivateRegionSecret{
							Clear: pointer.For(uuid.NewString()),
						},
					},
				}, nil, nil)

				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":              tftypes.String,
								"name":            tftypes.String,
								"organization_id": tftypes.String,
								"base_url":        tftypes.String,
								"secret":          tftypes.String,
							},
						}, map[string]tftypes.Value{
							"id":              tftypes.NewValue(tftypes.String, nil),
							"name":            tftypes.NewValue(tftypes.String, tc.name),
							"organization_id": tftypes.NewValue(tftypes.String, tc.organizationID),
							"base_url":        tftypes.NewValue(tftypes.String, ""),
							"secret":          tftypes.NewValue(tftypes.String, ""),
						}),
						Schema: resources.SchemaRegion,
					},
				}
				res := resource.CreateResponse{
					Diagnostics: []diag.Diagnostic{},
					State: tfsdk.State{
						Schema: resources.SchemaRegion,
					},
				}
				region.Create(ctx, req, &res)

				require.Len(t, res.Diagnostics, 0, "Expected no diagnostics on create")

				model := &resources.RegionModel{}
				res.State.Get(ctx, model)

				require.NotEmpty(t, model.ID.ValueString(), "Expected ID to be set")
				require.Equal(t, model.Name.ValueString(), tc.name, "Expected name to match")
				require.Equal(t, model.OrganizationID.ValueString(), tc.organizationID, "Expected organization ID to match")
				require.Equal(t, model.BaseURL.ValueString(), "https://example.com", "Expected base URL to be set")
				require.NotEmpty(t, model.Secret.ValueString(), "Expected secret to be set")

			})
		})
	}
}
