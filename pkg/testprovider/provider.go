package testprovider

import (
	"net/http"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/mock/gomock"
)

func NewCloudProvider(tracer trace.TracerProvider, logger logging.Logger, endpoint string, clientId string, clientSecret string, transport http.RoundTripper, sdkFactory pkg.CloudFactory, tokenFactory pkg.TokenProviderFactory) func() provider.Provider {
	return server.New(
		tracer,
		logger,
		endpoint,
		clientId,
		clientSecret,
		transport,
		sdkFactory,
		tokenFactory,
	)
}

type Mock struct {
	pkg.Creds
	*pkg.MockTokenProviderImpl
}

func NewMockTokenProvider(ctrl *gomock.Controller) (pkg.TokenProviderFactory, *Mock) {
	mock := &Mock{
		MockTokenProviderImpl: pkg.NewMockTokenProviderImpl(ctrl),
	}
	return func(transport http.RoundTripper, creds pkg.Creds) pkg.TokenProviderImpl {
		mock.Creds = creds
		return mock
	}, mock
}
