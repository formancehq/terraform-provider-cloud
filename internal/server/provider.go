package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/datasources"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	providerType = "formancecloud"
)

type FormanceCloudProviderModel struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Endpoint     types.String `tfsdk:"endpoint"`
}

type ProviderModelAdapter struct {
	m *FormanceCloudProviderModel
}

func NewProviderModelAdapter(m *FormanceCloudProviderModel) *ProviderModelAdapter {
	return &ProviderModelAdapter{
		m: m,
	}
}

func (f *ProviderModelAdapter) ClientId() string {
	return f.m.ClientId.ValueString()
}
func (f *ProviderModelAdapter) ClientSecret() string {
	return f.m.ClientSecret.ValueString()
}
func (f *ProviderModelAdapter) Endpoint() string {
	return f.m.Endpoint.ValueString()
}

func (f *ProviderModelAdapter) UserAgent() string {
	return fmt.Sprintf("terraform-provider-cloud/%s", internal.Version)
}

type FormanceCloudProvider struct {
	logger               logging.Logger
	transport            http.RoundTripper
	sdkFactory           pkg.SDKFactory
	tokenProviderFactory pkg.TokenProviderFactory

	Endpoint string

	ClientId     string
	ClientSecret string
}

var Schema = schema.Schema{
	Description: "The Formance Cloud provider allows you to manage your Formance Cloud resources using Terraform. It provides resources for managing organizations, stacks, regions, and stack modules.",
	Attributes: map[string]schema.Attribute{
		"client_secret": schema.StringAttribute{
			Description: "The client secret for authenticating with the Formance Cloud API. Can also be set via the FORMANCE_CLOUD_CLIENT_SECRET environment variable.",
			Optional:    true,
			Sensitive:   true,
		},
		"client_id": schema.StringAttribute{
			Description: "The client ID for authenticating with the Formance Cloud API. Can also be set via the FORMANCE_CLOUD_CLIENT_ID environment variable.",
			Optional:    true,
		},
		"endpoint": schema.StringAttribute{
			Description: "The endpoint URL for the Formance Cloud API. Defaults to the production endpoint. Can also be set via the FORMANCE_CLOUD_API_ENDPOINT environment variable.",
			Optional:    true,
		},
	},
}

// Metadata satisfies the provider.Provider interface for FormanceCloudProvider
func (p *FormanceCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = providerType
	resp.Version = internal.Version
}

// Schema satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = Schema
}

// Configure satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	p.logger.Debugf("Configuring cloud provider version %s", internal.Version)
	var data FormanceCloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.ClientId.ValueString() == "" {
		if p.ClientId != "" {
			data.ClientId = types.StringValue(p.ClientId)
		}
	}

	if data.ClientSecret.ValueString() == "" {
		if p.ClientSecret != "" {
			data.ClientSecret = types.StringValue(p.ClientSecret)
		}
	}

	if data.Endpoint.ValueString() == "" {
		data.Endpoint = types.StringValue(p.Endpoint)
	}

	creds := NewProviderModelAdapter(&data)
	tp := p.tokenProviderFactory(p.transport, creds)
	cli := p.sdkFactory(creds, pkg.NewTransport(p.transport, tp))

	resp.ResourceData = cli
	resp.DataSourceData = cli
}

// DataSources satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewOrganizations(p.logger.WithField("datasource", "organizations")),
		datasources.NewRegions(p.logger.WithField("datasource", "regions")),
		datasources.NewStacks(p.logger.WithField("datasource", "stacks")),
		datasources.NewRegionVersions(p.logger.WithField("datasource", "region_versions")),
	}
}

// Resources satisfies the provider.Provider interface for FormanceCloudProvider.
func (p *FormanceCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewStack(p.logger.WithField("resource", "stack")),
		resources.NewStackModule(p.logger.WithField("resource", "stack_module")),
		resources.NewStackMember(p.logger.WithField("resource", "stack_member")),
		resources.NewOrganizationMember(p.logger.WithField("resource", "organization_member")),
	}
}

func (p FormanceCloudProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{}
}

func (p FormanceCloudProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var data FormanceCloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if !data.ClientId.IsUnknown() {
		if p.ClientId != "" {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("client_id"),
				"Missing client_id Configuration",
				"While configuring the provider, the client_id was not found "+
					"However the FORMANCE_CLOUD_CLIENT_ID environment variable was set ",
			)
		} else {
			resp.Diagnostics.AddAttributeError(
				path.Root("client_id"),
				"Missing Client ID Configuration",
				"While configuring the provider, the client id was not found. "+
					"the FORMANCE_CLOUD_CLIENT_ID environment variable or provider "+
					"configuration block client_id attribute.",
			)
		}
	}

	if !data.ClientSecret.IsUnknown() {
		if p.ClientSecret != "" {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("client_secret"),
				"Missing client_secret Configuration",
				"While configuring the provider, the client_secret was not found in "+
					"however the environment variable FORMANCE_CLOUD_CLIENT_SECRET was found ",
			)
		} else {
			resp.Diagnostics.AddAttributeError(
				path.Root("client_secret"),
				"Missing API Token Configuration",
				"While configuring the provider, the API token was not found in "+
					"the FORMANCE_CLOUD_CLIENT_SECRET environment variable or provider "+
					"configuration block api_token attribute.",
			)
		}
	}

	if !data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("endpoint"),
			fmt.Sprintf("Missing Endpoint Configuration use %s", p.Endpoint),
			"While configuring the provider, the endpoint was not found "+
				"However the FORMANCE_CLOUD_API_ENDPOINT environment variable was set",
		)
	}
}

func New(
	logger logging.Logger,
	endpoint,
	clientId,
	clientSecret string,
	transport http.RoundTripper,
	sdkFactory pkg.SDKFactory,
	tokenFactory pkg.TokenProviderFactory,
) func() provider.Provider {
	return func() provider.Provider {
		return &FormanceCloudProvider{
			logger:               logger.WithField("provider", providerType),
			ClientId:             clientId,
			ClientSecret:         clientSecret,
			transport:            transport,
			Endpoint:             endpoint,
			sdkFactory:           sdkFactory,
			tokenProviderFactory: tokenFactory,
		}
	}
}

var _ provider.ProviderWithConfigValidators = &FormanceCloudProvider{}
var _ provider.ProviderWithValidateConfig = &FormanceCloudProvider{}
var _ provider.Provider = &FormanceCloudProvider{}
