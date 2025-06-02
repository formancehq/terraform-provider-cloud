package resources_test

import (
	"context"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider/internal/membership"
	"github.com/formancehq/terraform-provider/internal/membership/resources"
	"github.com/formancehq/terraform-provider/sdk"
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

func TestStackConfigure(t *testing.T) {
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
				providerData: membership.NewMockDefaultAPI(gomock.NewController(t)),
			},
		} {

			og := resources.NewStack(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

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

func TestStackMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		og := resources.NewStack(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

		res := resource.MetadataResponse{}

		og.Metadata(ctx, resource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_stack")
	})
}

func TestStackCreate(t *testing.T) {
	type testCase struct {
		name           string
		organizationID string
		regionID       string
		version        string
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
				r := resources.NewStack(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

				configureRes := resource.ConfigureResponse{
					Diagnostics: []diag.Diagnostic{},
				}
				ctrl := gomock.NewController(t)
				apiMock := membership.NewMockDefaultAPI(ctrl)
				r.Configure(ctx, resource.ConfigureRequest{
					ProviderData: apiMock,
				}, &configureRes)

				require.Empty(t, configureRes.Diagnostics, "Expected no diagnostics on configure")

				apiMock.EXPECT().CreateStack(gomock.Any(), tc.organizationID).Return(sdk.ApiCreateStackRequest{
					ApiService: apiMock,
				})

				md := map[string]string{
					"github.com/formancehq/terraform-provider/protected": "true",
				}
				expectedBody := sdk.CreateStackRequest{
					Name:     tc.name,
					Metadata: &md,
					RegionID: tc.regionID,
					Version:  &tc.version,
				}

				apiMock.EXPECT().CreateStackExecute(
					gomock.Cond(
						func(r sdk.ApiCreateStackRequest) bool {
							return cmp.Equal(r.GetCreateStackRequest(), &expectedBody, cmp.AllowUnexported(sdk.CreateStackRequest{}), cmpopts.IgnoreUnexported(sdk.CreateStackRequest{}))
						}),
				).Return(&sdk.CreateStackResponse{
					Data: &sdk.Stack{
						Id:             uuid.NewString(),
						Name:           tc.name,
						OrganizationId: tc.organizationID,
						RegionID:       tc.regionID,
						Version:        pointer.For(tc.version),
					},
				}, nil, nil)

				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":              tftypes.String,
								"name":            tftypes.String,
								"organization_id": tftypes.String,
								"region_id":       tftypes.String,
								"version":         tftypes.String,
								"force_destroy":   tftypes.Bool,
							},
						}, map[string]tftypes.Value{
							"id":              tftypes.NewValue(tftypes.String, nil),
							"name":            tftypes.NewValue(tftypes.String, tc.name),
							"organization_id": tftypes.NewValue(tftypes.String, tc.organizationID),
							"region_id":       tftypes.NewValue(tftypes.String, tc.regionID),
							"version":         tftypes.NewValue(tftypes.String, tc.version),
							"force_destroy":   tftypes.NewValue(tftypes.Bool, nil),
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
				r := resources.NewStack(logging.FromContext(ctx))().(resource.ResourceWithValidateConfig)

				res := resource.ValidateConfigResponse{
					Diagnostics: []diag.Diagnostic{},
				}

				r.ValidateConfig(ctx, resource.ValidateConfigRequest{
					Config: tfsdk.Config{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"name":            tftypes.String,
								"organization_id": tftypes.String,
								"region_id":       tftypes.String,
								"version":         tftypes.String,
								"id":              tftypes.String,
								"force_destroy":   tftypes.Bool,
							},
						}, map[string]tftypes.Value{
							"name":            tftypes.NewValue(tftypes.String, nil),
							"organization_id": tftypes.NewValue(tftypes.String, tc.organizationID),
							"region_id":       tftypes.NewValue(tftypes.String, tc.regionID),
							"version":         tftypes.NewValue(tftypes.String, nil),
							"id":              tftypes.NewValue(tftypes.String, nil),
							"force_destroy":   tftypes.NewValue(tftypes.Bool, nil),
						}),
						Schema: resources.SchemaStack,
					},
				}, &res)

				if tc.organizationID == nil || tc.regionID == nil {
					require.Len(t, res.Diagnostics, 2, "Expected one diagnostic for each missing field")

					require.Equal(t, res.Diagnostics[0].Summary(), "Invalid Organization ID")
					require.Equal(t, res.Diagnostics[1].Summary(), "Invalid Region ID")

				} else {
					require.Empty(t, res.Diagnostics, "Expected no diagnostics")
				}
			})
		})
	}
}
