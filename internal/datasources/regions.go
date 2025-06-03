package datasources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                   = &Region{}
	_ datasource.DataSourceWithConfigure      = &Region{}
	_ datasource.DataSourceWithValidateConfig = &Region{}
)

type Region struct {
	logger logging.Logger
	sdk    sdk.DefaultAPI
}

// ValidateConfig implements datasource.DataSourceWithValidateConfig.
func (r *Region) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, res *datasource.ValidateConfigResponse) {
	var config RegionModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.Name.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Name must be set.",
			"Region name cannot be empty.",
		)
	}

	if config.OrganizationID.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("organization_id"),
			"Organization ID must be set.",
			"Region organization ID cannot be null.",
		)
	}
}

var SchemaRegion = schema.Schema{
	Description: "Retrieves information about a specific region by name within an organization.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the region.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the region to retrieve.",
			Required:    true,
		},
		"organization_id": schema.StringAttribute{
			Description: "The organization ID where the region is located.",
			Required:    true,
		},
	},
}

// Configure implements datasource.DataSourceWithConfigure.
func (r *Region) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

	objs, res, err := r.sdk.ListRegions(ctx, data.OrganizationID.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}

	obj := collectionutils.First(objs.Data, func(o sdk.AnyRegion) bool {
		return o.Name == data.Name.ValueString()
	})
	if obj.Id == "" {
		resp.Diagnostics.AddError(
			"Region not found",
			fmt.Sprintf("No region found with name '%s' in organization '%s'", data.Name.ValueString(), data.OrganizationID.ValueString()),
		)
		return
	}

	data.ID = types.StringValue(obj.Id)
	data.Name = types.StringValue(obj.Name)
	data.OrganizationID = types.StringNull()
	if obj.OrganizationID != nil {
		data.OrganizationID = types.StringValue(*obj.OrganizationID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
