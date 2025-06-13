package datasources

import (
	"context"
	"fmt"
	"slices"
	"sort"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                   = &RegionVersions{}
	_ datasource.DataSourceWithConfigure      = &RegionVersions{}
	_ datasource.DataSourceWithValidateConfig = &RegionVersions{}
)

type RegionVersions struct {
	logger logging.Logger
	store  *pkg.Store
}

// ValidateConfig implements datasource.DataSourceWithValidateConfig.
func (r *RegionVersions) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, res *datasource.ValidateConfigResponse) {
	var config RegionVersionsModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.OrganizationID.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("organization_id"),
			"Organization ID must be set.",
			"RegionVersions organization ID cannot be null.",
		)
	}
}

var SchemaRegionVersions = schema.Schema{
	Description: "Retrieves the list of available Formance versions for a region. If id is specified, uses that region. Otherwise uses the first available region.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the region. If not specified, uses the first available region.",
			Optional:    true,
			Computed:    true,
		},
		"organization_id": schema.StringAttribute{
			Description: "The organization ID that owns the region.",
			Required:    true,
		},
		"versions": schema.ListNestedAttribute{
			Description: "The list of available Formance versions in the region.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "The version name (e.g., v1.0.0, v2.0.0).",
						Computed:    true,
					},
				},
			},
		},
	},
}

// Configure implements datasource.DataSourceWithConfigure.
func (r *RegionVersions) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

type RegionVersionsModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Versions       []Version    `tfsdk:"versions"`
}

type Version struct {
	Name types.String `tfsdk:"name"`
}

func NewRegionVersions(logger logging.Logger) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &RegionVersions{
			logger: logger,
		}
	}
}

func (r *RegionVersions) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_region_versions"
}

func (r *RegionVersions) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SchemaRegionVersions
}

func (r *RegionVersions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RegionVersionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var regionID string

	if !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != "" {
		// If ID is specified, use it
		regionID = data.ID.ValueString()
	} else {
		// If ID is not specified, list regions and use the first one
		regions, res, err := r.store.GetSDK().ListRegions(ctx, data.OrganizationID.ValueString()).Execute()
		if err != nil {
			pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
			return
		}

		if len(regions.Data) == 0 {
			resp.Diagnostics.AddError(
				"No regions found",
				fmt.Sprintf("No regions found in organization '%s'", data.OrganizationID.ValueString()),
			)
			return
		}

		// Sort regions deterministically by ID to ensure consistent selection
		sort.Slice(regions.Data, func(i, j int) bool {
			return regions.Data[i].Id < regions.Data[j].Id
		})

		// Use the first region after sorting
		regionID = regions.Data[0].Id
		data.ID = types.StringValue(regionID)
	}

	obj, res, err := r.store.GetSDK().GetRegionVersions(ctx, data.OrganizationID.ValueString(), regionID).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
		return
	}

	versions := make([]Version, len(obj.Data))
	for i, v := range obj.Data {
		versions[i] = Version{
			Name: types.StringValue(v.Name),
		}
	}
	slices.SortFunc(versions, func(a, b Version) int {
		if a.Name.ValueString() < b.Name.ValueString() {
			return -1
		}
		if a.Name.ValueString() > b.Name.ValueString() {
			return 1
		}
		return 0
	})

	data.Versions = versions

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
