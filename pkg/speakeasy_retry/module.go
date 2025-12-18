package speakeasyretry

import (
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/retry"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
)

var (
	RetryFlag                = "retry-enabled"
	RetryInitialIntervalFlag = "retry-initial-interval"
	RetryMaxIntervalFlag     = "retry-max-interval"
	RetryMaxElapsedTimeFlag  = "retry-max-elapsed-time"
	RetryExponentFlag        = "retry-exponent"
)

const (
	DefaultRetryInitialInterval = 1000
	DefaultRetryMaxInterval     = 3000
	DefaultRetryMaxElapsedTime  = 10000
	DefaultRetryExponent        = 2.0
)

func AddFlags(flags *pflag.FlagSet) {
	flags.Bool(RetryFlag, true, "Enable SDK retry")
	flags.Int(RetryInitialIntervalFlag, DefaultRetryInitialInterval, "Initial interval for retry backoff strategy in milliseconds")
	flags.Int(RetryMaxIntervalFlag, DefaultRetryMaxInterval, "Max interval for retry backoff strategy in milliseconds")
	flags.Int(RetryMaxElapsedTimeFlag, DefaultRetryMaxElapsedTime, "Max elapsed time for retry backoff strategy in milliseconds")
	flags.Float64(RetryExponentFlag, DefaultRetryExponent, "Exponent for retry backoff strategy")
}

func NewModule(flags *pflag.FlagSet) fx.Option {
	initialInterval, _ := flags.GetInt(RetryInitialIntervalFlag)
	maxInterval, _ := flags.GetInt(RetryMaxIntervalFlag)
	maxElapsedTime, _ := flags.GetInt(RetryMaxElapsedTimeFlag)
	exponent, _ := flags.GetFloat64(RetryExponentFlag)

	if enabled, _ := flags.GetBool(RetryFlag); !enabled {
		return fx.Options(
			fx.Provide(func() *retry.Config {
				return nil
			}),
		)
	}
	return fx.Options(
		fx.Supply(&retry.Config{
			Strategy: "backoff",
			Backoff: &retry.BackoffStrategy{
				InitialInterval: initialInterval,
				MaxInterval:     maxInterval,
				Exponent:        exponent,
				MaxElapsedTime:  maxElapsedTime,
			},
			RetryConnectionErrors: true,
		}),
	)
}
