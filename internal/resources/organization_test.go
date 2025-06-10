package resources_test

import (
	"context"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
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
				providerData: pkg.NewMockDefaultAPI(gomock.NewController(t)),
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
