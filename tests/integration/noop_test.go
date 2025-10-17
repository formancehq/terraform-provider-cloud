package integration_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/server"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	type testCase struct {
		clientId     string
		clientSecret string
		endpoint     string

		expectedError string
	}

	for _, tc := range []testCase{
		{
			clientId:     fmt.Sprintf("organization_%s", uuid.NewString()[:8]),
			clientSecret: uuid.NewString(),
			endpoint:     "https://app.formance.cloud/api",
		},
		{
			clientId:      uuid.NewString(),
			clientSecret:  uuid.NewString(),
			endpoint:      "https://app.formance.cloud/api",
			expectedError: "Invalid client_id",
		},
		{
			clientId:      uuid.NewString(),
			endpoint:      "https://app.formance.cloud/api",
			expectedError: "Missing client_secret",
		},
		{
			clientId:      uuid.NewString(),
			expectedError: "Missing endpoint",
		},
	} {
		t.Run(fmt.Sprintf("%s %+v", t.Name(), tc), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			cloudSdk := pkg.NewMockCloudSDK(ctrl)
			tokenProvider := pkg.NewMockTokenProviderImpl(ctrl)

			cloudProvider := server.NewProvider(
				noop.NewTracerProvider(),
				logging.Testing().WithField("provider", "cloud_noop"),
				server.FormanceCloudEndpoint(tc.endpoint),
				server.FormanceCloudClientId(tc.clientId),
				server.FormanceCloudClientSecret(tc.clientSecret),
				transport,
				func(creds pkg.Creds, transport http.RoundTripper) pkg.CloudSDK {
					return cloudSdk
				},
				func(transport http.RoundTripper, creds pkg.Creds) pkg.TokenProviderImpl {
					return tokenProvider
				},
			)

			noopStep := resource.TestStep{
				Config: `
					provider "cloud" {}

					resource "cloud_noop" "default" {}

				`,
			}

			if tc.expectedError != "" {
				noopStep.ExpectError = regexp.MustCompile(tc.expectedError)
			}

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"cloud": providerserver.NewProtocol6WithError(cloudProvider()),
				},
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version0_15_0),
				},
				Steps: []resource.TestStep{
					noopStep,
				},
			})
		})
	}

}
