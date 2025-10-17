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
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"go.uber.org/mock/gomock"
)

func TestStacksConfigure(t *testing.T) {
	type testCase struct {
		providerData func(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) any
		expectedErr  error
	}

	for _, tc := range []testCase{
		{
			providerData: func(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) any {
				return nil
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
		ctx := logging.TestingContext()
		og := datasources.NewStacks()().(datasource.DataSourceWithConfigure)

		res := datasource.ConfigureResponse{
			Diagnostics: []diag.Diagnostic{},
		}
		ctrl := gomock.NewController(t)
		tp := pkg.NewMockTokenProviderImpl(ctrl)
		apiMock := pkg.NewMockCloudSDK(ctrl)
		data := tc.providerData(apiMock, tp)

		if tc.expectedErr == nil && data != nil {
			tp.EXPECT().IntrospectToken(gomock.Any()).Return(oidc.IntrospectionResponse{
				Claims: map[string]interface{}{
					"organization_id": "test-org-id",
				},
			}, nil).AnyTimes()
		}
		og.Configure(ctx, datasource.ConfigureRequest{
			ProviderData: data,
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
	og := datasources.NewStacks()().(datasource.DataSourceWithConfigure)

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

			og := datasources.NewStacks()().(datasource.DataSourceWithConfigValidators)

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
