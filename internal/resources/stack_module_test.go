package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/pkg"

	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStackModuleConfigure(t *testing.T) {
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
				providerData: internal.NewStore(pkg.NewMockDefaultAPI(gomock.NewController(t)), fmt.Sprintf("organization_%s", uuid.NewString())),
			},
		} {

			og := resources.NewStackModule(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

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

func TestStackModuleMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		og := resources.NewStackModule(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

		res := resource.MetadataResponse{}

		og.Metadata(ctx, resource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_stack_module")
	})
}

func TestStackModuleValidateConfig(t *testing.T) {
	type testCase struct {
		name           *string
		organizationID *string
		stackID        *string
	}

	for _, tc := range []testCase{
		{},
		{
			name:           pointer.For(uuid.NewString()),
			organizationID: pointer.For(uuid.NewString()),
			stackID:        pointer.For(uuid.NewString()),
		},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			test(t, func(ctx context.Context) {
				og := resources.NewStackModule(logging.FromContext(ctx))().(resource.ResourceWithValidateConfig)

				res := resource.ValidateConfigResponse{
					Diagnostics: []diag.Diagnostic{},
				}

				og.ValidateConfig(ctx, resource.ValidateConfigRequest{
					Config: tfsdk.Config{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: getSchemaTypes(resources.SchemaStackModule),
						}, map[string]tftypes.Value{
							"name":     tftypes.NewValue(tftypes.String, tc.name),
							"stack_id": tftypes.NewValue(tftypes.String, tc.stackID),
						}),
						Schema: resources.SchemaStackModule,
					},
				}, &res)

				if tc.name == nil || tc.organizationID == nil || tc.stackID == nil {
					require.Len(t, res.Diagnostics, 2, "Expected one diagnostic for missing name")
					require.Equal(t, res.Diagnostics[0].Summary(), "Invalid Name")
					require.Equal(t, res.Diagnostics[1].Summary(), "Invalid Stack ID")
				} else {
					require.Empty(t, res.Diagnostics, "Expected no diagnostics")
				}
			})
		})
	}
}
