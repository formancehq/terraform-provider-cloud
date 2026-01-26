package integration_test

import (
	"net/http"

	"github.com/formancehq/terraform-provider-cloud/pkg"
)

func NewCloudSdkMockT(mock *pkg.MockCloudSDK) func(endpoint string, transport http.RoundTripper) pkg.CloudSDK {
	return func(endpoint string, transport http.RoundTripper) pkg.CloudSDK {
		return mock
	}
}

func NewCloudTokenProviderMockT(mock *pkg.MockTokenProviderImpl) func(transport http.RoundTripper, creds pkg.Creds, scopes []string) pkg.TokenProviderImpl {
	return func(transport http.RoundTripper, creds pkg.Creds, scopes []string) pkg.TokenProviderImpl {
		return mock
	}
}
