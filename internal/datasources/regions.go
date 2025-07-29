package datasources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
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
	logger logging.Logger
	store  *internal.Store
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
	ctx = logging.ContextWithLogger(ctx, r.logger.WithField("func", "region_read"))
	r.logger.Debug("Reading region")
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
	var obj sdk.AnyRegion
	switch {
	case !data.ID.IsNull():
		objs, res, err := r.store.GetSDK().GetRegion(ctx, organizationId, data.ID.ValueString())
		if err != nil {
			pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
			return
		}
		obj = objs.Data
	case !data.Name.IsNull():
		objs, res, err := r.store.GetSDK().ListRegions(ctx, organizationId)
		if err != nil {
			pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
			return
		}
		obj = collectionutils.First(objs.Data, func(o sdk.AnyRegion) bool {
			return o.Name == data.Name.ValueString()
		})
		if obj.Id == "" {
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

	data.ID = types.StringValue(obj.Id)
	data.Name = types.StringValue(obj.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
