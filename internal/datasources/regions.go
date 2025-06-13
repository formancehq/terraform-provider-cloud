package datasources

import (
	"context"
	"fmt"
	"sort"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                   = &Region{}
	_ datasource.DataSourceWithConfigure      = &Region{}
	_ datasource.DataSourceWithValidateConfig = &Region{}
)

type Region struct {
	logger logging.Logger
	store  *pkg.Store
}

// ValidateConfig implements datasource.DataSourceWithValidateConfig.
func (r *Region) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, res *datasource.ValidateConfigResponse) {
	var config RegionModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	// Organization ID is now optional - will use Store if not provided
}

var SchemaRegion = schema.Schema{
	Description: "Retrieves information about regions within an organization. If name is specified, returns a specific region by name.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the region.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the region to retrieve. If not specified, returns the first available region.",
			Optional:    true,
			Computed:    true,
		},
		"organization_id": schema.StringAttribute{
			Description: "The organization ID where the region is located. If not specified, uses the current organization.",
			Optional:    true,
			Computed:    true,
		},
	},
}

// Configure implements datasource.DataSourceWithConfigure.
func (r *Region) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
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

	r.store = store
}

type RegionModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	OrganizationID types.String `tfsdk:"organization_id"`
}

func NewRegions(logger logging.Logger) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &Region{
			logger: logger,
		}
	}
}

func (r *Region) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regions"
}

func (r *Region) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SchemaRegion
}

func (r *Region) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RegionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use organization ID from config or fall back to the store
	orgID := data.OrganizationID.ValueString()
	if orgID == "" {
		orgID = r.store.GetOrganizationID()
		if orgID == "" {
			// Try to fetch and set the current organization
			var err error
			orgID, err = r.store.FetchAndSetCurrentOrganization(ctx)
			if err != nil {
				pkg.HandleSDKError(ctx, err, nil, &resp.Diagnostics)
				return
			}
			if orgID == "" {
				resp.Diagnostics.AddError(
					"No organization found",
					"Unable to determine organization ID. Please specify organization_id or ensure the user has access to at least one organization.",
				)
				return
			}
		}
	}

	objs, res, err := r.store.GetSDK().ListRegions(ctx, orgID).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
		return
	}

	var obj sdk.AnyRegion
	
	if !data.Name.IsNull() && !data.Name.IsUnknown() && data.Name.ValueString() != "" {
		// If name is specified, find the specific region
		obj = collectionutils.First(objs.Data, func(o sdk.AnyRegion) bool {
			return o.Name == data.Name.ValueString()
		})
		if obj.Id == "" {
			resp.Diagnostics.AddError(
				"Region not found",
				fmt.Sprintf("No region found with name '%s' in organization '%s'", data.Name.ValueString(), orgID),
			)
			return
		}
	} else {
		// If name is not specified, sort and return the first available region
		if len(objs.Data) == 0 {
			resp.Diagnostics.AddError(
				"No regions found",
				fmt.Sprintf("No regions found in organization '%s'", orgID),
			)
			return
		}
		
		// Sort regions deterministically by ID to ensure consistent selection
		sort.Slice(objs.Data, func(i, j int) bool {
			return objs.Data[i].Id < objs.Data[j].Id
		})
		
		obj = objs.Data[0]
	}

	data.ID = types.StringValue(obj.Id)
	data.Name = types.StringValue(obj.Name)
	// Set the organization ID that was actually used
	data.OrganizationID = types.StringValue(orgID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
