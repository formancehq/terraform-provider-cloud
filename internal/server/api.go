package server

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

type ProviderFactory func() provider.Provider

type API struct {
	provider ProviderFactory
}

func (a *API) Run(ctx context.Context, debug bool) error {
	opts := providerserver.ServeOpts{
		Address: fmt.Sprintf("%s/%s", "registry.terraform.io", internal.Repository),
		Debug:   debug,
	}

	err := providerserver.Serve(ctx, a.provider, opts)
	if err != nil {
		logging.FromContext(ctx).Errorf("failed to start server: %v", err)
		return err
	}

	return nil
}
func NewAPI(p ProviderFactory) *API {
	return &API{
		provider: p,
	}
}
