package server

import (
	"context"

	"github.com/formancehq/go-libs/v3/service"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
)

const (
	FormanceCloudClientSecretKey = "formance-cloud-client-secret"
	FormanceCloudClientIdKey     = "formance-cloud-client-id"
	FormanceCloudEndpointKey     = "formance-cloud-api-endpoint"
)

func AddFlags(flagset *pflag.FlagSet) {
	flagset.String(FormanceCloudClientSecretKey, "", "User Client Secret for Formance Cloud")
	flagset.String(FormanceCloudClientIdKey, "", "User ID for Formance Cloud")
	flagset.String(FormanceCloudEndpointKey, "https://app.formance.cloud/api", "Endpoint for Formance Cloud")
}

func NewModule(ctx context.Context, flagset *pflag.FlagSet) fx.Option {
	clientId, _ := flagset.GetString(FormanceCloudClientIdKey)
	clientSecret, _ := flagset.GetString(FormanceCloudClientSecretKey)
	endpoint, _ := flagset.GetString(FormanceCloudEndpointKey)
	opts := fx.Options(
		fx.Supply(FormanceCloudClientId(clientId)),
		fx.Supply(FormanceCloudClientSecret(clientSecret)),
		fx.Supply(FormanceCloudEndpoint(endpoint)),
		fx.Provide(
			NewAPI,
		),
		fx.Invoke(func(lc fx.Lifecycle, server *API, shutdowner fx.Shutdowner) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					debug, _ := flagset.GetBool(service.DebugFlag)
					go func() {
						if err := server.Run(ctx, debug); err != nil {
							if err := shutdowner.Shutdown(); err != nil {
								panic(err)
							}
						}
					}()
					return nil
				},
			})
		}),
	)

	return opts
}
