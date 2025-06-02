package e2e_test

import (
	"os"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider/internal/membership"
	"github.com/formancehq/terraform-provider/internal/providers"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

var Provider func() provider.Provider

func TestMain(m *testing.M) {
	endpoint := os.Getenv("FORMANCE_CLOUD_API_ENDPOINT")
	clientID := os.Getenv("FORMANCE_CLOUD_CLIENT_ID")
	clientSecret := os.Getenv("FORMANCE_CLOUD_CLIENT_SECRET")
	Provider = providers.New(logging.Testing(), "develop", endpoint, clientID, clientSecret, membership.NewSDK)
	code := m.Run()

	os.Exit(code)
}
