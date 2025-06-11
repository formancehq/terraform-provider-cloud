package pkg

import (
	"net/http"

	"github.com/formancehq/terraform-provider-cloud/sdk"
	gomock "go.uber.org/mock/gomock"
)

type Creds interface {
	ClientId() string
	ClientSecret() string
	Endpoint() string
	UserAgent() string
}

//go:generate rm -rf ../sdk
//go:generate openapi-generator-cli generate -i ./openapi.yaml -g go -o ../sdk --git-user-id=formancehq --git-repo-id=terraform-provider-cloud -p packageVersion=latest -p isGoSubmodule=true -p packageName=sdk -p disallowAdditionalPropertiesIfNotPresent=false -p generateInterfaces=true -t ../openapi-templates/go
//go:generate rm -rf ../sdk/test
//go:generate rm -rf ../sdk/docs
func NewSDK(creds Creds, transport http.RoundTripper) (sdk.DefaultAPI, TokenProviderImpl) {
	tp := NewTokenProvider(transport, creds)

	client := http.Client{
		Transport: newTransport(transport, tp),
	}

	sdk := &SDK{
		APIClient: sdk.NewAPIClient(&sdk.Configuration{
			HTTPClient: &client,
			UserAgent:  creds.UserAgent(),
			Servers: sdk.ServerConfigurations{
				{
					URL:         creds.Endpoint(),
					Description: "Membership API",
				},
			},
		}),
	}
	return sdk.DefaultAPI, &tp
}

type Mocks struct {
	Api   *MockDefaultAPI
	Creds Creds

	TokenProvider *MockTokenProviderImpl
}

func NewMockSDK(ctrl *gomock.Controller) (SDKFactory, *Mocks) {
	mockSDK := NewMockDefaultAPI(ctrl)
	mockTokenProvider := NewMockTokenProviderImpl(ctrl)
	mocks := &Mocks{
		Api:           mockSDK,
		TokenProvider: mockTokenProvider,
	}
	return func(creds Creds, transport http.RoundTripper) (sdk.DefaultAPI, TokenProviderImpl) {
		mocks.Creds = creds
		return mockSDK, mockTokenProvider
	}, mocks
}

type SDKFactory func(creds Creds, transport http.RoundTripper) (sdk.DefaultAPI, TokenProviderImpl)

//go:generate mockgen -source=../sdk/api_default.go -destination=sdk_generated.go -package=pkg . DefaultAPI
type SDK struct {
	*sdk.APIClient
}

type RoundTripperFn func(r *http.Request) (*http.Response, error)

func (fn RoundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
func newTransport(rt http.RoundTripper, tp TokenProviderImpl) RoundTripperFn {
	return func(r *http.Request) (*http.Response, error) {
		token, err := tp.RefreshToken(r.Context())
		if err != nil {
			return nil, err
		}
		r.Header.Set("Authorization", "Bearer "+token.AccessToken)
		return rt.RoundTrip(r)
	}
}
