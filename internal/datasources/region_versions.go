package datasources

import (
	"context"
	"fmt"
	"slices"
	"sort"

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
	store *internal.Store
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

func NewRegionVersions() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &RegionVersions{}
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
	organizationId, err := r.store.GetOrganizationID(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}
	var regionID string

	if !data.ID.IsNull() {
		regionID = data.ID.ValueString()
	} else {
		operation, err := r.store.GetSDK().ListRegions(ctx, organizationId)
		if err != nil {
			pkg.HandleSDKError(ctx, err, &resp.Diagnostics)
			return
		}

		if len(operation.ListRegionsResponse.Data) == 0 {
			resp.Diagnostics.AddError(
				"No regions found",
				fmt.Sprintf("No regions found in organization '%s'", organizationId),
			)
			return
		}

		sort.Slice(operation.ListRegionsResponse.Data, func(i, j int) bool {
			return operation.ListRegionsResponse.Data[i].ID < operation.ListRegionsResponse.Data[j].ID
		})

		regionID = operation.ListRegionsResponse.Data[0].ID
		data.ID = types.StringValue(regionID)
	}

	operation, err := r.store.GetSDK().GetRegionVersions(ctx, organizationId, regionID)
	if err != nil {
		pkg.HandleSDKError(ctx, err, &resp.Diagnostics)
		return
	}

	versions := make([]Version, len(operation.GetRegionVersionsResponse.Data))
	for i, v := range operation.GetRegionVersionsResponse.Data {
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
