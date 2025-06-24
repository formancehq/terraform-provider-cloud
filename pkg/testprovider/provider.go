package testprovider

import (
	"net/http"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func NewCloudProvider(logger logging.Logger, endpoint string, clientId string, clientSecret string, transport http.RoundTripper, sdkFactory pkg.CloudFactory, tokenFactory pkg.TokenProviderFactory) func() provider.Provider {
	return server.New(
		logger,
		endpoint,
		clientId,
		clientSecret,
		transport,
		sdkFactory,
		tokenFactory,
	)
}
