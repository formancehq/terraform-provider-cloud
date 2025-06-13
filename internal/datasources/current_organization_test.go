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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCurrentOrganizationConfigure(t *testing.T) {
	test(t, func(ctx context.Context) {

		type testCase struct {
			providerData     any
			expectedErrTitle string
		}

		for _, tc := range []testCase{
			{
				expectedErrTitle: "Provider Not Configured",
			},
			{
				providerData:     "something",
				expectedErrTitle: resources.ErrProviderDataNotSet.Error(),
			},
			{
				providerData: pkg.NewMockDefaultAPI(gomock.NewController(t)),
			},
		} {

			co := datasources.NewCurrentOrganization(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

			res := datasource.ConfigureResponse{
				Diagnostics: []diag.Diagnostic{},
			}

			co.Configure(ctx, datasource.ConfigureRequest{
				ProviderData: tc.providerData,
			}, &res)

			if tc.expectedErrTitle != "" {
				require.Len(t, res.Diagnostics, 1, "Expected one diagnostic")
				require.Equal(t, res.Diagnostics[0].Summary(), tc.expectedErrTitle)
			} else {
				require.Empty(t, res.Diagnostics, "Expected no diagnostics")
			}

		}
	})
}

func TestCurrentOrganizationMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		co := datasources.NewCurrentOrganization(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

		res := datasource.MetadataResponse{}

		co.Metadata(ctx, datasource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_current_organization")
	})
}

func TestCurrentOrganizationValidateConfig(t *testing.T) {
	test(t, func(ctx context.Context) {
		co := datasources.NewCurrentOrganization(logging.FromContext(ctx))().(datasource.DataSourceWithValidateConfig)
		res := datasource.ValidateConfigResponse{
			Diagnostics: []diag.Diagnostic{},
		}
		co.ValidateConfig(ctx, datasource.ValidateConfigRequest{}, &res)

		// Current organization has no required config, so should have no diagnostics
		require.Empty(t, res.Diagnostics, "Expected no diagnostics")
	})
}