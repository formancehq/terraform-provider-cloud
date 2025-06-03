package testprovider

import (
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func NewCloudProvider(logger logging.Logger, version string, endpoint string, clientId string, clientSecret string, sdkFactory pkg.SDKFactory) func() provider.Provider {
	return server.New(
		logger,
		version,
		endpoint,
		clientId,
		clientSecret,
		sdkFactory,
	)
}
