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
	_ datasource.DataSource                   = &Stack{}
	_ datasource.DataSourceWithConfigure      = &Stack{}
	_ datasource.DataSourceWithValidateConfig = &Stack{}
)

type Stack struct {
	logger logging.Logger
	sdk    sdk.DefaultAPI
}

// ValidateConfig implements datasource.DataSourceWithValidateConfig.
func (s *Stack) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, res *datasource.ValidateConfigResponse) {
	var config StackModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.ID.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("id"),
			"ID must be set.",
			"ID cannot be empty.",
		)
	}

	if config.OrganizationID.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("organization_id"),
			"Organization ID must be set.",
			"Organization ID cannot be null.",
		)
	}
}

var SchemaStack = schema.Schema{
	Description: "Retrieves information about a specific Formance Cloud stack by ID.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the stack to retrieve.",
			Required:    true,
		},
		"organization_id": schema.StringAttribute{
			Description: "The organization ID that owns the stack.",
			Required:    true,
		},
	},
}

// Configure implements datasource.DataSourceWithConfigure.
func (s *Stack) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

	s.sdk = sdk
}

type StackModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
}

func NewStacks(logger logging.Logger) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &Stack{
			logger: logger,
		}
	}
}

func (s *Stack) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stacks"
}

func (s *Stack) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SchemaStack
}

func (s *Stack) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StackModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	obj, res, err := s.sdk.GetStack(ctx, data.OrganizationID.ValueString(), data.ID.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}

	if obj == nil {
		resp.Diagnostics.AddError("Unable to read stack", "Stack not found")
		return
	}
	data.ID = types.StringValue(obj.Data.Id)
	data.OrganizationID = types.StringValue(obj.Data.OrganizationId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
