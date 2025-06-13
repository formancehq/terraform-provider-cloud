package datasources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                   = &CurrentOrganization{}
	_ datasource.DataSourceWithConfigure      = &CurrentOrganization{}
	_ datasource.DataSourceWithValidateConfig = &CurrentOrganization{}
)

var SchemaCurrentOrganization = schema.Schema{
	Description: "Retrieves information about the current/first organization associated with the authenticated user.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the organization.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the organization.",
			Computed:    true,
		},
		"owner_id": schema.StringAttribute{
			Description: "The ID of the organization owner.",
			Computed:    true,
		},
		"domain": schema.StringAttribute{
			Description: "The domain of the organization.",
			Computed:    true,
		},
	},
}

type CurrentOrganization struct {
	logger logging.Logger
	store  *pkg.Store
}

// ValidateConfig implements datasource.DataSourceWithValidateConfig.
func (c *CurrentOrganization) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, res *datasource.ValidateConfigResponse) {
	// No configuration to validate for current organization
}

// Configure implements datasource.DataSourceWithConfigure.
func (c *CurrentOrganization) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		res.Diagnostics.AddError(
			"Provider Not Configured",
			"The provider has not been configured before use, please ensure the provider configuration is set.",
		)
		return
	}

	store, ok := req.ProviderData.(*pkg.Store)
	if !ok {
		res.Diagnostics.AddError(
			resources.ErrProviderDataNotSet.Error(),
			fmt.Sprintf("Expected *pkg.Store, got: %T", req.ProviderData),
		)
		return
	}

	c.store = store
}

type CurrentOrganizationModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	OwnerID types.String `tfsdk:"owner_id"`
	Domain  types.String `tfsdk:"domain"`
}

func NewCurrentOrganization(logger logging.Logger) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &CurrentOrganization{
			logger: logger,
		}
	}
}

func (c *CurrentOrganization) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_current_organization"
}

func (c *CurrentOrganization) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SchemaCurrentOrganization
}

func (c *CurrentOrganization) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CurrentOrganizationModel

	ctx = logging.ContextWithLogger(ctx, c.logger.WithField("func", "current_organization_read"))
	logging.FromContext(ctx).Debugf("Reading current organization")

	// Try to get cached organization ID first
	orgID := c.store.GetOrganizationID()
	
	// If no cached ID, fetch from API
	if orgID == "" {
		var err error
		orgID, err = c.store.FetchAndSetCurrentOrganization(ctx)
		if err != nil {
			pkg.HandleSDKError(ctx, err, nil, &resp.Diagnostics)
			return
		}
		
		if orgID == "" {
			resp.Diagnostics.AddError(
				"No organizations found",
				"The authenticated user does not have access to any organizations.",
			)
			return
		}
	}
	
	// Fetch all organizations and find the one with matching ID
	orgsResp, res, err := c.store.GetSDK().ListOrganizationsExpanded(ctx).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
		return
	}
	
	// Find the organization with the cached ID
	var org *sdk.OrganizationExpanded
	for _, o := range orgsResp.Data {
		if o.Id == orgID {
			org = &o
			break
		}
	}
	
	if org == nil {
		resp.Diagnostics.AddError(
			"Organization not found",
			fmt.Sprintf("Organization with ID '%s' not found", orgID),
		)
		return
	}

	data.ID = types.StringValue(org.Id)
	data.Name = types.StringValue(org.Name)
	data.OwnerID = types.StringValue(org.OwnerId)
	
	if org.Domain != nil {
		data.Domain = types.StringValue(*org.Domain)
	} else {
		data.Domain = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}