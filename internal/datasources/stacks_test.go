package datasources_test

import (
	"fmt"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/datasources"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStacksConfigure(t *testing.T) {
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
			providerData: internal.NewStore(pkg.NewMockCloudSDK(gomock.NewController(t)), fmt.Sprintf("organization_%s", uuid.NewString())),
		},
	} {
		ctx := logging.TestingContext()
		og := datasources.NewStacks(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

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

}

func TestStacksMetadata(t *testing.T) {
	ctx := logging.TestingContext()
	og := datasources.NewStacks(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

	res := datasource.MetadataResponse{}

	og.Metadata(ctx, datasource.MetadataRequest{
		ProviderTypeName: "test",
	}, &res)

	require.Contains(t, res.TypeName, "_stacks")

}

func TestStacksConfigValidator(t *testing.T) {
	type testCase struct {
		id          *string
		name        *string
		expectedErr error
	}

	for i, tc := range []testCase{
		{
			expectedErr: fmt.Errorf("Missing Attribute Configuration"),
		},
		{
			id:   pointer.For(uuid.NewString()),
			name: pointer.For(uuid.NewString()),
		},
		{
			id: pointer.For(uuid.NewString()),
		},
		{
			name: pointer.For(uuid.NewString()),
		},
	} {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			t.Parallel()
			ctx := logging.TestingContext()

			og := datasources.NewStacks(logging.FromContext(ctx))().(datasource.DataSourceWithConfigValidators)

			schemaStackType := getSchemaTypes(datasources.SchemaStack)
			validators := og.ConfigValidators(ctx)
			for _, validator := range validators {
				res := datasource.ValidateConfigResponse{
					Diagnostics: []diag.Diagnostic{},
				}
				validator.ValidateDataSource(ctx, datasource.ValidateConfigRequest{
					Config: tfsdk.Config{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: schemaStackType,
							OptionalAttributes: map[string]struct{}{
								"region_id": {},
								"state":     {},
								"status":    {},
							},
						}, map[string]tftypes.Value{
							"name": tftypes.NewValue(tftypes.String, tc.name),
							"id":   tftypes.NewValue(tftypes.String, tc.id),
						}),
						Schema: datasources.SchemaStack,
					},
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
}
