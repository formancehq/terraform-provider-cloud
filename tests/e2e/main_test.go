package e2e_test

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/httpclient"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/otlp"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

var (
	Provider func() provider.Provider
	RegionName string
	OrganizationId string
)

func TestMain(m *testing.M) {
	endpoint := os.Getenv("FORMANCE_CLOUD_API_ENDPOINT")
	clientID := os.Getenv("FORMANCE_CLOUD_CLIENT_ID")
	clientSecret := os.Getenv("FORMANCE_CLOUD_CLIENT_SECRET")

	flag.Parse()

	var transport http.RoundTripper
	if testing.Verbose() {
		transport = httpclient.NewDebugHTTPTransport(
			otlp.NewRoundTripper(http.DefaultTransport, true),
		)
	} else {
		transport = otlp.NewRoundTripper(http.DefaultTransport, false)
	}

	Provider = server.New(logging.Testing(), endpoint, clientID, clientSecret, transport, pkg.NewSDK, pkg.NewTokenProvider)

	// Setup non destroyable resources
	RegionName = os.Getenv("FORMANCE_CLOUD_REGION_NAME")
	OrganizationId = os.Getenv("FORMANCE_CLOUD_ORGANIZATION_ID")

	s := []string{RegionName, OrganizationId, endpoint, clientID, clientSecret}
	s = collectionutils.Filter(s, func(s string) bool {
		return s == ""
	})

	if len(s) > 0 {
		fmt.Printf("Missing environment variables: %+v\n", s)
		panic("You must set the following environment variables")
	}

	code := m.Run()

	os.Exit(code)
}
