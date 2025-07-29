package datasources

import (
	"context"
	"fmt"
	"slices"
	"sort"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &RegionVersions{}
	_ datasource.DataSourceWithConfigure = &RegionVersions{}
)

type RegionVersions struct {
	logger logging.Logger
	store  *internal.Store
}

var SchemaRegionVersions = schema.Schema{
	Description: "Retrieves the list of available Formance versions for a region. If id is specified, uses that region. Otherwise, uses the first available region sorted deterministically by ID.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the region. If not specified, uses the first available region sorted deterministically by ID.",
			Optional:    true,
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

type RegionVersionsModel struct {
	ID       types.String `tfsdk:"id"`
	Versions []Version    `tfsdk:"versions"`
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
	ctx = logging.ContextWithLogger(ctx, r.logger.WithField("func", "region_versions_read"))
	r.logger.Debug("Reading region versions")
	var data RegionVersionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var regionID string

	if !data.ID.IsNull() {
		regionID = data.ID.ValueString()
	} else {
		regions, res, err := r.store.GetSDK().ListRegions(ctx, r.store.GetOrganizationID(ctx))
		if err != nil {
			pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
			return
		}

		if len(regions.Data) == 0 {
			resp.Diagnostics.AddError(
				"No regions found",
				fmt.Sprintf("No regions found in organization '%s'", r.store.GetOrganizationID(ctx)),
			)
			return
		}

		sort.Slice(regions.Data, func(i, j int) bool {
			return regions.Data[i].Id < regions.Data[j].Id
		})

		regionID = regions.Data[0].Id
		data.ID = types.StringValue(regionID)
	}

	obj, res, err := r.store.GetSDK().GetRegionVersions(ctx, r.store.GetOrganizationID(ctx), regionID)
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
