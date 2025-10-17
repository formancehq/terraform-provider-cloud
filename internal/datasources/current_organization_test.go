package datasources_test

import (
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/datasources"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"go.uber.org/mock/gomock"
)

func TestCurrentOrganizationConfigure(t *testing.T) {

	type testCase struct {
		providerData  func(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) any
		expectedError error
	}

	for _, tc := range []testCase{
		{
			providerData: func(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) any {
				return "something"
			},
			expectedError: resources.ErrProviderDataNotSet,
		},
		{
			providerData: func(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) any {
				return internal.NewStore(sdkClient, tp)
			},
		},
	} {
		ctx := logging.TestingContext()
		co := datasources.NewCurrentOrganization()().(datasource.DataSourceWithConfigure)

		res := datasource.ConfigureResponse{
			Diagnostics: []diag.Diagnostic{},
		}

		ctrl := gomock.NewController(t)
		tp := pkg.NewMockTokenProviderImpl(ctrl)
		apiMock := pkg.NewMockCloudSDK(ctrl)
		data := tc.providerData(apiMock, tp)
		if tc.expectedError == nil && data != nil {
			tp.EXPECT().IntrospectToken(gomock.Any()).Return(oidc.IntrospectionResponse{
				Claims: map[string]interface{}{
					"organization_id": "test-org-id",
				},
			}, nil).AnyTimes()
		}

		co.Configure(ctx, datasource.ConfigureRequest{
			ProviderData: data,
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
	co := datasources.NewCurrentOrganization()().(datasource.DataSourceWithConfigure)

	res := datasource.MetadataResponse{}

	co.Metadata(ctx, datasource.MetadataRequest{
		ProviderTypeName: "test",
	}, &res)

	require.Contains(t, res.TypeName, "_current_organization")

}
