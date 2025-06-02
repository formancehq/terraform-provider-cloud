package datasources

import (
	"context"
	"fmt"
	"slices"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider/internal/membership"
	"github.com/formancehq/terraform-provider/internal/membership/resources"
	"github.com/formancehq/terraform-provider/sdk"
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
	sdk    sdk.DefaultAPI
}

// ValidateConfig implements datasource.DataSourceWithValidateConfig.
func (r *RegionVersions) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, res *datasource.ValidateConfigResponse) {
	var config RegionVersionsModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}
	if config.ID.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("id"),
			"ID must be set.",
			"RegionVersions ID cannot be empty.",
		)
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
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Required: true,
		},
		"organization_id": schema.StringAttribute{
			Required: true,
		},
		"versions": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed: true,
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

	sdk, ok := req.ProviderData.(sdk.DefaultAPI)
	if !ok {
		res.Diagnostics.AddError(
			resources.ErrProviderDataNotSet.Error(),
			fmt.Sprintf("Expected *FormanceCloudProviderModel, got: %T", req.ProviderData),
		)
		return
	}

	r.sdk = sdk
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

	obj, res, err := r.sdk.GetRegionVersions(ctx, data.OrganizationID.ValueString(), data.ID.ValueString()).Execute()
	if err != nil {
		membership.HandleSDKError(ctx, res, &resp.Diagnostics)
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
