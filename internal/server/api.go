package server

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

type FormanceCloudEndpoint string
type FormanceCloudClientSecret string
type FormanceCloudClientId string

type API struct {
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

	err := providerserver.Serve(ctx, New(a.logger, internal.Version, string(a.Endpoint), string(a.ClientId), string(a.ClientSecret), pkg.NewSDK), opts)
	if err != nil {
		logging.FromContext(ctx).Errorf("failed to start server: %v", err)
		return err
	}

	return nil
}
func NewAPI(endpoint FormanceCloudEndpoint, clientSecret FormanceCloudClientSecret, clientId FormanceCloudClientId, logger logging.Logger) *API {
	return &API{
		logger:       logger,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     endpoint,
	}
}
