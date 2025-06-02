package datasources_test

import (
	"context"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/datasources"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestOrganizationsConfigure(t *testing.T) {
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

			og := datasources.NewOrganizations(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

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

func TestOrganizationsMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		og := datasources.NewOrganizations(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

		res := datasource.MetadataResponse{}

		og.Metadata(ctx, datasource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_organizations")
	})
}

func TestOrganizationValidateConfig(t *testing.T) {
	type testCase struct {
		id *string
	}

	for _, tc := range []testCase{} {
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := logging.TestingContext()

			og := datasources.NewOrganizations(logging.FromContext(ctx))().(datasource.DataSourceWithValidateConfig)
			res := datasource.ValidateConfigResponse{
				Diagnostics: []diag.Diagnostic{},
			}
			og.ValidateConfig(ctx, datasource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"id": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"id": tftypes.NewValue(tftypes.String, tc.id),
					}),
					Schema: datasources.SchemaOrganization,
				},
			}, &res)

			if tc.id != nil {
				require.Len(t, res.Diagnostics, 1, "Expected one diagnostic")
				require.Equal(t, res.Diagnostics[0].Summary(), "ID must be set.")
			} else {
				require.Empty(t, res.Diagnostics, "Expected no diagnostics")
			}

		})
	}
}
