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

func TestOrganizationConfigure(t *testing.T) {
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

			og := resources.NewOrganization(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

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

func TestOrganizationMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		og := resources.NewOrganization(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

		res := resource.MetadataResponse{}

		og.Metadata(ctx, resource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_organization")
	})
}

func TestOrganizationValidateConfig(t *testing.T) {
	type testCase struct {
		name *string
	}

	for _, tc := range []testCase{
		{
			name: pointer.For(uuid.NewString()),
		},
		{},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			test(t, func(ctx context.Context) {
				og := resources.NewOrganization(logging.FromContext(ctx))().(resource.ResourceWithValidateConfig)

				res := resource.ValidateConfigResponse{
					Diagnostics: []diag.Diagnostic{},
				}

				og.ValidateConfig(ctx, resource.ValidateConfigRequest{
					Config: tfsdk.Config{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":                          tftypes.String,
								"name":                        tftypes.String,
								"domain":                      tftypes.String,
								"default_organization_access": tftypes.String,
								"default_stack_access":        tftypes.String,
							},
						}, map[string]tftypes.Value{
							"id":                          tftypes.NewValue(tftypes.String, nil),
							"name":                        tftypes.NewValue(tftypes.String, tc.name),
							"domain":                      tftypes.NewValue(tftypes.String, nil),
							"default_organization_access": tftypes.NewValue(tftypes.String, nil),
							"default_stack_access":        tftypes.NewValue(tftypes.String, nil),
						}),
						Schema: resources.SchemaOrganization,
					},
				}, &res)

				if tc.name == nil {
					require.Len(t, res.Diagnostics, 1, "Expected one diagnostic on validate config")
				} else {
					require.Empty(t, res.Diagnostics, "Expected no diagnostics on validate config")
				}
			})
		})
	}
}

func TestOrganizationCreate(t *testing.T) {
	type testCase struct {
		name                      string
		domain                    string
		defaultOrganizationAccess string
		defaultStackAccess        string
	}

	for _, tc := range []testCase{
		{
			name:                      "test-org",
			domain:                    "test-domain",
			defaultOrganizationAccess: "ADMIN",
			defaultStackAccess:        "ADMIN",
		},
		{
			name:                      "test-org",
			defaultOrganizationAccess: "ADMIN",
			defaultStackAccess:        "ADMIN",
		},
		{
			name: "test-org",
		},
		{},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			test(t, func(ctx context.Context) {
				og := resources.NewOrganization(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

				configureRes := resource.ConfigureResponse{
					Diagnostics: []diag.Diagnostic{},
				}
				ctrl := gomock.NewController(t)
				apiMock := membership.NewMockDefaultAPI(ctrl)
				og.Configure(ctx, resource.ConfigureRequest{
					ProviderData: apiMock,
				}, &configureRes)

				require.Empty(t, configureRes.Diagnostics, "Expected no diagnostics on configure")

				res := resource.CreateResponse{
					Diagnostics: []diag.Diagnostic{},
					State: tfsdk.State{
						Schema: resources.SchemaOrganization,
					},
				}

				values := map[string]tftypes.Value{
					"id":                          tftypes.NewValue(tftypes.String, nil),
					"name":                        tftypes.NewValue(tftypes.String, tc.name),
					"domain":                      tftypes.NewValue(tftypes.String, tc.domain),
					"default_organization_access": tftypes.NewValue(tftypes.String, tc.defaultOrganizationAccess),
					"default_stack_access":        tftypes.NewValue(tftypes.String, tc.defaultStackAccess),
				}
				req := resource.CreateRequest{
					Plan: tfsdk.Plan{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id":                          tftypes.String,
								"name":                        tftypes.String,
								"domain":                      tftypes.String,
								"default_organization_access": tftypes.String,
								"default_stack_access":        tftypes.String,
							},
						}, values),
						Schema: resources.SchemaOrganization,
					},
				}

				apiMock.EXPECT().CreateOrganization(gomock.Any()).Return(sdk.ApiCreateOrganizationRequest{
					ApiService: apiMock,
				})

				expectedBody := sdk.OrganizationData{
					Name: tc.name,
				}

				if tc.domain != "" {
					expectedBody.Domain = pointer.For(tc.domain)
				}
				if tc.defaultOrganizationAccess != "" {
					expectedBody.DefaultOrganizationAccess = pointer.For(sdk.Role(tc.defaultOrganizationAccess))
				}
				if tc.defaultStackAccess != "" {
					expectedBody.DefaultStackAccess = pointer.For(sdk.Role(tc.defaultStackAccess))
				}
				apiMock.EXPECT().CreateOrganizationExecute(
					gomock.Cond(
						func(r sdk.ApiCreateOrganizationRequest) bool {
							return cmp.Equal(r.GetBody(), &expectedBody, cmp.AllowUnexported(sdk.OrganizationData{}), cmpopts.IgnoreUnexported(sdk.OrganizationData{}))
						}),
				).Return(&sdk.CreateOrganizationResponse{
					Data: &sdk.OrganizationExpanded{
						Id:                        uuid.NewString(),
						Name:                      tc.name,
						Domain:                    pointer.For(tc.domain),
						DefaultOrganizationAccess: pointer.For(sdk.Role(tc.defaultOrganizationAccess)),
						DefaultStackAccess:        pointer.For(sdk.Role(tc.defaultStackAccess)),
					},
				}, nil, nil)
				og.Create(ctx, req, &res)

				require.Len(t, res.Diagnostics, 0, "Expected no diagnostics on create")

				model := &resources.OrganizationModel{}
				res.State.Get(ctx, model)

				require.Equal(t, model.Name.ValueString(), tc.name, "Expected name to match")
				require.Equal(t, model.Domain.ValueString(), tc.domain, "Expected domain to match")
				require.Equal(t, model.DefaultOrganizationAccess.ValueString(), tc.defaultOrganizationAccess, "Expected default organization access to match")
				require.Equal(t, model.DefaultStackAccess.ValueString(), tc.defaultStackAccess, "Expected default stack access to match")
				require.NotEmpty(t, model.ID.ValueString(), "Expected ID to be set")
			})
		})
	}
}
