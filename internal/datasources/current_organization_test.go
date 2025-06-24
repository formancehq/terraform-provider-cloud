package datasources_test

import (
	"fmt"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/datasources"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCurrentOrganizationConfigure(t *testing.T) {

	type testCase struct {
		providerData  any
		expectedError error
	}

	for _, tc := range []testCase{
		{
			providerData:  "something",
			expectedError: resources.ErrProviderDataNotSet,
		},
		{
			providerData: internal.NewStore(pkg.NewMockCloudSDK(gomock.NewController(t)), fmt.Sprintf("organization_%s", uuid.NewString())),
		},
	} {
		ctx := logging.TestingContext()
		co := datasources.NewCurrentOrganization(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

		res := datasource.ConfigureResponse{
			Diagnostics: []diag.Diagnostic{},
		}

		co.Configure(ctx, datasource.ConfigureRequest{
			ProviderData: tc.providerData,
		}, &res)

		if tc.expectedError != nil {
			require.Len(t, res.Diagnostics, 1, "Expected one diagnostic")
			require.Equal(t, res.Diagnostics[0].Summary(), tc.expectedError.Error())
		} else {
			require.Empty(t, res.Diagnostics, "Expected no diagnostics")
		}

	}

}

func TestCurrentOrganizationMetadata(t *testing.T) {
	ctx := logging.TestingContext()
	co := datasources.NewCurrentOrganization(logging.FromContext(ctx))().(datasource.DataSourceWithConfigure)

	res := datasource.MetadataResponse{}

	co.Metadata(ctx, datasource.MetadataRequest{
		ProviderTypeName: "test",
	}, &res)

	require.Contains(t, res.TypeName, "_current_organization")

}
