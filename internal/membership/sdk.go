package membership

import (
	"net/http"

	"github.com/formancehq/go-libs/v3/httpclient"
	"github.com/formancehq/go-libs/v3/otlp"
	"github.com/formancehq/terraform-provider/internal"
	"github.com/formancehq/terraform-provider/sdk"
	"go.uber.org/mock/gomock"
)

//go:generate rm -rf ./sdk
//go:generate openapi-generator-cli generate -i ./openapi.yaml -g go -o ./sdk --git-user-id=formancehq --git-repo-id=terraform-provider -p packageVersion=latest -p isGoSubmodule=true -p packageName=sdk -p disallowAdditionalPropertiesIfNotPresent=false -p generateInterfaces=true -t ../../openapi-templates/go
//go:generate rm -rf ./sdk/test
func NewSDK(version string, creds *internal.FormanceCloudProviderModel) (sdk.DefaultAPI, TokenProviderImpl) {
	client := http.Client{
		Transport: otlp.NewRoundTripper(httpclient.NewDebugHTTPTransport(http.DefaultTransport), true),
	}

	tp := NewTokenProvider(&http.Client{
		Transport: otlp.NewRoundTripper(httpclient.NewDebugHTTPTransport(http.DefaultTransport), true),
	}, creds)
	client.Transport = newTransport(client.Transport, tp.cloud)
	sdk := &SDK{
		APIClient: sdk.NewAPIClient(&sdk.Configuration{
			HTTPClient: &client,
			UserAgent:  "terraform-provider/" + version,
			Servers: sdk.ServerConfigurations{
				{
					URL:         creds.Endpoint.ValueString(),
					Description: "Membership API",
				},
			},
		}),
	}
	return sdk.DefaultAPI, &tp
}

type Mocks struct {
	Api     *MockDefaultAPI
	Creds   *internal.FormanceCloudProviderModel
	Version string

	TokenProvider *MockTokenProviderImpl
}

func NewMockSDK(ctrl *gomock.Controller) (SDKFactory, *Mocks) {
	mockSDK := NewMockDefaultAPI(ctrl)
	mockTokenProvider := NewMockTokenProviderImpl(ctrl)
	mocks := &Mocks{
		Api:           mockSDK,
		TokenProvider: mockTokenProvider,
	}
	return func(version string, creds *internal.FormanceCloudProviderModel) (sdk.DefaultAPI, TokenProviderImpl) {
		mocks.Version = version
		mocks.Creds = creds
		return mockSDK, mockTokenProvider
	}, mocks
}

type SDKFactory func(version string, creds *internal.FormanceCloudProviderModel) (sdk.DefaultAPI, TokenProviderImpl)

//go:generate mockgen -source=./sdk/api_default.go -destination=sdk_generated.go -package=membership . DefaultAPI
type SDK struct {
	*sdk.APIClient
}

type RoundTripperFn func(r *http.Request) (*http.Response, error)

func (fn RoundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
func newTransport(rt http.RoundTripper, m *internal.TokenInfo) RoundTripperFn {
	return func(r *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()

		r.Header.Set("Authorization", "Bearer "+m.AccessToken)
		return rt.RoundTrip(r)
	}
}
