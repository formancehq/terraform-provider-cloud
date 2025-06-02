package providers_test

import (
	"fmt"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider/internal/membership"
	"github.com/formancehq/terraform-provider/internal/providers"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestProviderMetadata(t *testing.T) {
	p := providers.New(logging.Testing(), "develop", "https://app.formance.cloud/api", "client_id", "client_secret", membership.NewSDK)()

	res := provider.MetadataResponse{}
	p.Metadata(logging.TestingContext(), provider.MetadataRequest{}, &res)

	require.Equal(t, res.TypeName, "formancecloud")
	require.Equal(t, res.Version, "develop")
}

func TestProviderConfigure(t *testing.T) {
	type testCase struct {
		ClientId     string
		ClientSecret string
		Endpoint     string
	}

	for _, tc := range []testCase{
		{
			ClientId:     uuid.NewString(),
			ClientSecret: uuid.NewString(),
			Endpoint:     uuid.NewString(),
		},
		{},
	} {
		t.Run(fmt.Sprintf("%s clientId %t clientSecret %t endpoint %t", t.Name(), tc.ClientId != "", tc.ClientSecret != "", tc.Endpoint != ""), func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			sdkFactory, mocks := membership.NewMockSDK(ctrl)
			p := providers.New(logging.Testing(), "develop", "https://app.formance.cloud/api", "client_id", "client_secret", sdkFactory)()

			mocks.TokenProvider.EXPECT().RefreshToken(gomock.Any()).Return(nil)

			res := provider.ConfigureResponse{
				Diagnostics: []diag.Diagnostic{},
			}

			p.Configure(logging.TestingContext(), provider.ConfigureRequest{
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"client_id":     tftypes.String,
							"client_secret": tftypes.String,
							"endpoint":      tftypes.String,
						},
					}, map[string]tftypes.Value{
						"client_id":     tftypes.NewValue(tftypes.String, tc.ClientId),
						"client_secret": tftypes.NewValue(tftypes.String, tc.ClientSecret),
						"endpoint":      tftypes.NewValue(tftypes.String, tc.Endpoint),
					}),
					Schema: providers.Schema,
				},
			}, &res)

			require.Equal(t, mocks.Version, "develop")
			if tc.ClientId == "" {
				require.Equal(t, mocks.Creds.ClientId.ValueString(), "client_id")
			}

			if tc.ClientSecret == "" {
				require.Equal(t, mocks.Creds.ClientSecret.ValueString(), "client_secret")
			}

			if tc.Endpoint == "" {
				require.Equal(t, mocks.Creds.Endpoint.ValueString(), "https://app.formance.cloud/api")
			}

			require.Len(t, res.Diagnostics, 0)
			require.NotNil(t, res.ResourceData)
			require.NotNil(t, res.DataSourceData)
		})
	}

}
