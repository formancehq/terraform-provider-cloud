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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                   = &Organization{}
	_ datasource.DataSourceWithConfigure      = &Organization{}
	_ datasource.DataSourceWithValidateConfig = &Organization{}
)

var SchemaOrganization = schema.Schema{
	Description: "Retrieves information about a specific Formance Cloud organization by ID.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the organization to retrieve.",
			Required:    true,
		},
	},
}

type Organization struct {
	logger logging.Logger
	sdk    sdk.DefaultAPI
}

// ValidateConfig implements datasource.DataSourceWithValidateConfig.
func (o *Organization) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, res *datasource.ValidateConfigResponse) {
	var config OrganizationModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.ID.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("id"),
			"ID must be set.",
			"Organization ID cannot be empty.",
		)
	}
}

// Configure implements datasource.DataSourceWithConfigure.
func (o *Organization) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

	o.sdk = sdk
}

type OrganizationModel struct {
	ID types.String `tfsdk:"id"`
}

func NewOrganizations(logger logging.Logger) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &Organization{
			logger: logger,
		}
	}
}

func (o *Organization) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (o *Organization) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SchemaOrganization
}

func (o *Organization) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	ctx = logging.ContextWithLogger(ctx, o.logger.WithField("func", "organization_read"))
	logging.FromContext(ctx).Debugf("Reading organization")

	if resp.Diagnostics.HasError() {
		return
	}

	obj, res, err := o.sdk.ReadOrganization(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
		return
	}

	data.ID = types.StringValue(obj.Data.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
