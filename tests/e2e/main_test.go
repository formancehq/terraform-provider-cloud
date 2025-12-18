package e2e_test

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/formancehq/go-libs/v3/httpclient"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/otlp"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var (
	Provider       func() provider.Provider
	RegionName     string
	OrganizationId string
)

func TestMain(m *testing.M) {
	endpoint := os.Getenv("FORMANCE_CLOUD_API_ENDPOINT")
	clientID := os.Getenv("FORMANCE_CLOUD_CLIENT_ID")
	clientSecret := os.Getenv("FORMANCE_CLOUD_CLIENT_SECRET")

	flag.Parse()

	var (
		transport http.RoundTripper
		tp        trace.TracerProvider
	)
	if testing.Verbose() {
		p := sdktrace.NewTracerProvider()
		tp = p
		transport = httpclient.NewDebugHTTPTransport(
			otlp.NewRoundTripper(http.DefaultTransport, true),
		)
		defer func() {
			if err := p.ForceFlush(context.Background()); err != nil {
				panic(err)
			}
		}()
	} else {
		tp = noop.NewTracerProvider()
		transport = otlp.NewRoundTripper(http.DefaultTransport, false)
	}
	Provider = server.New(
		tp,
		logging.Testing(),
		endpoint,
		clientID,
		clientSecret,
		transport,
		pkg.NewCloudSDK,
		pkg.NewTokenProvider,
	)

	// Setup non destroyable resources
	RegionName = os.Getenv("FORMANCE_CLOUD_REGION_NAME")
	OrganizationId = os.Getenv("FORMANCE_CLOUD_ORGANIZATION_ID")

	// Check only required variables
	requiredVars := map[string]string{
		"FORMANCE_CLOUD_API_ENDPOINT":    endpoint,
		"FORMANCE_CLOUD_CLIENT_ID":       clientID,
		"FORMANCE_CLOUD_CLIENT_SECRET":   clientSecret,
		"FORMANCE_CLOUD_REGION_NAME":     RegionName,
		"FORMANCE_CLOUD_ORGANIZATION_ID": OrganizationId,
	}

	var missingVars []string
	for name, value := range requiredVars {
		if value == "" {
			missingVars = append(missingVars, name)
		}
	}

	if len(missingVars) > 0 {
		fmt.Printf("Missing required environment variables: %+v\n", missingVars)
		panic("You must set the required environment variables")
	}

	code := m.Run()

	os.Exit(code)
}
