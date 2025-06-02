package server

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider/internal"
	"github.com/formancehq/terraform-provider/internal/membership"
	"github.com/formancehq/terraform-provider/internal/providers"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

type FormanceCloudEndpoint string
type FormanceCloudClientSecret string
type FormanceCloudClientId string

type API struct {
	internal.AppInfo

	logger       logging.Logger
	ClientId     FormanceCloudClientId
	ClientSecret FormanceCloudClientSecret
	Endpoint     FormanceCloudEndpoint
}

func (a *API) Run(ctx context.Context, debug bool) error {
	opts := providerserver.ServeOpts{
		Address: fmt.Sprintf("%s/%s", "registry.terraform.io", internal.Repository),
		Debug:   debug,
	}

	err := providerserver.Serve(ctx, providers.New(a.logger, a.Version, string(a.Endpoint), string(a.ClientId), string(a.ClientSecret), membership.NewSDK), opts)
	if err != nil {
		logging.FromContext(ctx).Errorf("failed to start server: %v", err)
		return err
	}

	return nil
}
func New(info internal.AppInfo, endpoint FormanceCloudEndpoint, clientSecret FormanceCloudClientSecret, clientId FormanceCloudClientId, logger logging.Logger) *API {
	return &API{
		logger:       logger,
		AppInfo:      info,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     endpoint,
	}
}
