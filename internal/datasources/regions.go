package datasources

import (
	"context"
	"fmt"

	"github.com/formancehq/formance-sdk-cloud-go/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                     = &Region{}
	_ datasource.DataSourceWithConfigure        = &Region{}
	_ datasource.DataSourceWithConfigValidators = &Region{}
)

type Region struct {
	store *internal.Store
}

var SchemaRegion = schema.Schema{
	Description: "Retrieves information about regions within an organization. If name is specified, returns a specific region by name. Otherwise, returns the first available region sorted deterministically by ID.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the region.",
			Optional:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the region to retrieve. If not specified, returns the first available region sorted deterministically by ID.",
			Optional:    true,
		},
	},
}

// ConfigValidators implements datasource.DataSourceWithConfigValidators.
func (r *Region) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

// Configure implements datasource.DataSourceWithConfigure.
func (r *Region) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	store, ok := req.ProviderData.(*internal.Store)
	if !ok {
		res.Diagnostics.AddError(
			resources.ErrProviderDataNotSet.Error(),
			fmt.Sprintf("Expected *internal.Store, got: %T", req.ProviderData),
		)
		return
	}

	r.store = store
}

type RegionModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewRegions() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &Region{}
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
	organizationId, err := r.store.GetOrganizationID(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}
	var region shared.AnyRegion
	switch {
	case !data.ID.IsNull():
		operation, err := r.store.GetSDK().GetRegion(ctx, organizationId, data.ID.ValueString())
		if err != nil {
			pkg.HandleSDKError(ctx, err, &resp.Diagnostics)
			return
		}
		region = operation.GetRegionResponse.Data
	case !data.Name.IsNull():
		operation, err := r.store.GetSDK().ListRegions(ctx, organizationId)
		if err != nil {
			pkg.HandleSDKError(ctx, err, &resp.Diagnostics)
			return
		}
		region = collectionutils.First(operation.ListRegionsResponse.Data, func(o shared.AnyRegion) bool {
			return o.Name == data.Name.ValueString()
		})
		if region.ID == "" {
			resp.Diagnostics.AddError(
				"Region not found",
				fmt.Sprintf("No region found with name '%s' in organization '%s'", data.Name.ValueString(), organizationId),
			)
			return
		}

	default:
		resp.Diagnostics.AddError(
			"Region ID or Name required",
			"Either 'id' or 'name' must be specified to retrieve a region.",
		)
		return
	}

	data.ID = types.StringValue(region.ID)
	data.Name = types.StringValue(region.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
