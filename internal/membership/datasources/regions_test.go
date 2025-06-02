package datasources_test

import (
	"context"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider/internal/membership"
	"github.com/formancehq/terraform-provider/internal/membership/datasources"
	"github.com/formancehq/terraform-provider/internal/membership/resources"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegionsConfigure(t *testing.T) {
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

			og := datasources.NewRegions(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

			res := datasource.ConfigureResponse{
				Diagnostics: []diag.Diagnostic{},
			}

			og.Configure(ctx, datasource.ConfigureRequest{
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

func TestRegionsMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		og := datasources.NewRegions(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

		res := datasource.MetadataResponse{}

		og.Metadata(ctx, datasource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_regions")
	})
}

func TestRegionsValidateConfig(t *testing.T) {
	type testCase struct {
		name           *string
		organizationID *string
	}

	for _, tc := range []testCase{
		{},
		{
			name:           pointer.For(uuid.NewString()),
			organizationID: pointer.For(uuid.NewString()),
		},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := logging.TestingContext()

			r := datasources.NewRegions(logging.FromContext(ctx))().(datasource.DataSourceWithValidateConfig)
			res := datasource.ValidateConfigResponse{
				Diagnostics: []diag.Diagnostic{},
			}
			r.ValidateConfig(ctx, datasource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"name":            tftypes.String,
							"organization_id": tftypes.String,
							"id":              tftypes.String,
						},
					}, map[string]tftypes.Value{
						"name":            tftypes.NewValue(tftypes.String, tc.name),
						"organization_id": tftypes.NewValue(tftypes.String, tc.organizationID),
						"id":              tftypes.NewValue(tftypes.String, nil),
					}),
					Schema: datasources.SchemaRegion,
				},
			}, &res)
			if tc.name == nil && tc.organizationID == nil {
				require.Len(t, res.Diagnostics, 2, "Expected diagnostics")
				require.Equal(t, res.Diagnostics[0].Summary(), "Name must be set.")
				require.Equal(t, res.Diagnostics[1].Summary(), "Organization ID must be set.")
			} else {
				require.Empty(t, res.Diagnostics, "Expected no diagnostics")
			}
		})
	}
}
