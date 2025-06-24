package server_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestProviderMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	tokenFactory, _ := pkg.NewMockTokenProvider(ctrl)
	p := server.New(logging.Testing(), "https://app.formance.cloud/api", "client_id", "client_secret", http.DefaultTransport, pkg.NewCloudSDK(), tokenFactory)()

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
			tokenFactory, mockTp := pkg.NewMockTokenProvider(ctrl)
			p := server.New(logging.Testing(), "https://app.formance.cloud/api", "client_id", "client_secret", http.DefaultTransport, pkg.NewCloudSDK(), tokenFactory)()

			res := provider.ConfigureResponse{
				Diagnostics: []diag.Diagnostic{},
			}

			p.Configure(logging.TestingContext(), provider.ConfigureRequest{
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: getSchemaTypes(server.Schema),
					}, map[string]tftypes.Value{
						"client_id":     tftypes.NewValue(tftypes.String, tc.ClientId),
						"client_secret": tftypes.NewValue(tftypes.String, tc.ClientSecret),
						"endpoint":      tftypes.NewValue(tftypes.String, tc.Endpoint),
					}),
					Schema: server.Schema,
				},
			}, &res)

			if tc.ClientId == "" {
				require.Equal(t, mockTp.ClientId(), "client_id")
			}

			if tc.ClientSecret == "" {
				require.Equal(t, mockTp.ClientSecret(), "client_secret")
			}

			if tc.Endpoint == "" {
				require.Equal(t, mockTp.Endpoint(), "https://app.formance.cloud/api")
			}

			require.Len(t, res.Diagnostics, 0)
			require.NotNil(t, res.ResourceData)
			require.NotNil(t, res.DataSourceData)
		})
	}

}
