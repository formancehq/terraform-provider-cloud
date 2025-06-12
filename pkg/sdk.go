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
func NewSDK(creds Creds, transport http.RoundTripper) sdk.DefaultAPI {
	sdk := &SDK{
		APIClient: sdk.NewAPIClient(&sdk.Configuration{
			HTTPClient: &http.Client{
				Transport: transport,
			},
			UserAgent: creds.UserAgent(),
			Servers: sdk.ServerConfigurations{
				{
					URL:         creds.Endpoint(),
					Description: "Membership API",
				},
			},
		}),
	}
	return sdk.DefaultAPI
}

func NewSdkFactory() SDKFactory {
	return func(creds Creds, transport http.RoundTripper) sdk.DefaultAPI {
		return NewSDK(creds, transport)
	}
}

func NewMockSDK(ctrl *gomock.Controller) (SDKFactory, *MockDefaultAPI) {
	mock := NewMockDefaultAPI(ctrl)
	return func(creds Creds, transport http.RoundTripper) sdk.DefaultAPI {
		return mock
	}, mock
}

type SDKFactory func(creds Creds, transport http.RoundTripper) sdk.DefaultAPI

//go:generate mockgen -source=../sdk/api_default.go -destination=sdk_generated.go -package=pkg . DefaultAPI
type SDK struct {
	*sdk.APIClient
}

type RoundTripperFn func(r *http.Request) (*http.Response, error)

func (fn RoundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
func NewTransport(rt http.RoundTripper, tp TokenProviderImpl) RoundTripperFn {
	return func(r *http.Request) (*http.Response, error) {
		token, err := tp.RefreshToken(r.Context())
		if err != nil {
			return nil, err
		}
		r.Header.Set("Authorization", "Bearer "+token.AccessToken)
		return rt.RoundTrip(r)
	}
}
