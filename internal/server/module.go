package server

import (
	"context"
	"net/http"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/otlp"
	"github.com/formancehq/go-libs/v3/service"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/retry"
	speakeasyretry "github.com/formancehq/terraform-provider-cloud/pkg/speakeasy_retry"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/trace"
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
	speakeasyretry.AddFlags(flagset)
}

type FormanceCloudEndpoint string
type FormanceCloudClientSecret string
type FormanceCloudClientId string

func NewProvider(
	tracer trace.TracerProvider,
	logger logging.Logger,
	endpoint FormanceCloudEndpoint,
	clientId FormanceCloudClientId,
	clientSecret FormanceCloudClientSecret,
	transport http.RoundTripper,
	sdkFactory pkg.CloudFactory,
	tokenFactory pkg.TokenProviderFactory,
) ProviderFactory {
	return New(
		tracer,
		logger,
		string(endpoint),
		string(clientId),
		string(clientSecret),
		transport,
		sdkFactory,
		tokenFactory,
	)
}

func NewModule(ctx context.Context, flagset *pflag.FlagSet) fx.Option {
	clientId, _ := flagset.GetString(FormanceCloudClientIdKey)
	clientSecret, _ := flagset.GetString(FormanceCloudClientSecretKey)
	endpoint, _ := flagset.GetString(FormanceCloudEndpointKey)
	debug, _ := flagset.GetBool(service.DebugFlag)
	transport := otlp.NewRoundTripper(http.DefaultTransport, debug)
	return fx.Options(
		fx.Supply(FormanceCloudClientId(clientId)),
		fx.Supply(FormanceCloudClientSecret(clientSecret)),
		fx.Supply(FormanceCloudEndpoint(endpoint)),
		fx.Supply(fx.Annotate(transport, fx.As(new(http.RoundTripper)))),
		speakeasyretry.NewModule(flagset),
		fx.Provide(func() pkg.TokenProviderFactory {
			return pkg.NewTokenProvider
		}),
		fx.Provide(func(retry *retry.Config) pkg.CloudFactory {
			opts := []membershipclient.SDKOption{}
			if retry != nil {
				opts = append(opts, membershipclient.WithRetryConfig(*retry))
			}
			return pkg.NewCloudSDK(
				opts...,
			)
		}),
		fx.Provide(NewProvider),
		fx.Provide(NewAPI),
		fx.Invoke(func(lc fx.Lifecycle, server *API, shutdowner fx.Shutdowner) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
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
}
