package datasources

import (
	"context"
	"fmt"
	"sort"

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
	store  *pkg.Store
}

// ValidateConfig implements datasource.DataSourceWithValidateConfig.
func (s *Stack) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, res *datasource.ValidateConfigResponse) {
	var config StackModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
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
	Description: "Retrieves information about a Formance Cloud stack. If id is specified, returns a specific stack by ID. Otherwise, returns the first available stack sorted alphabetically by name for predictable behavior.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the stack. If not specified, returns the first available stack sorted alphabetically by name.",
			Optional:    true,
			Computed:    true,
		},
		"organization_id": schema.StringAttribute{
			Description: "The organization ID that owns the stack.",
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the stack.",
			Computed:    true,
		},
		"region_id": schema.StringAttribute{
			Description: "The region ID where the stack is installed.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The current status of the stack.",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "The current state of the stack.",
			Computed:    true,
		},
	},
}

// Configure implements datasource.DataSourceWithConfigure.
func (s *Stack) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

	s.store = store
}

type StackModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	RegionID       types.String `tfsdk:"region_id"`
	Status         types.String `tfsdk:"status"`
	State          types.String `tfsdk:"state"`
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
	if resp.Diagnostics.HasError() {
		return
	}

	var stack sdk.Stack

	if !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != "" {
		// If ID is specified, get the specific stack
		obj, res, err := s.store.GetSDK().GetStack(ctx, data.OrganizationID.ValueString(), data.ID.ValueString()).Execute()
		if err != nil {
			pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
			return
		}

		if obj == nil || obj.Data == nil {
			resp.Diagnostics.AddError("Unable to read stack", "Stack not found")
			return
		}
		stack = *obj.Data
	} else {
		// If ID is not specified, list all stacks and return the first one (sorted deterministically)
		listResp, res, err := s.store.GetSDK().ListStacks(ctx, data.OrganizationID.ValueString()).Execute()
		if err != nil {
			pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
			return
		}

		if len(listResp.Data) == 0 {
			resp.Diagnostics.AddError(
				"No stacks found",
				fmt.Sprintf("No stacks found in organization '%s'", data.OrganizationID.ValueString()),
			)
			return
		}

		// Sort stacks deterministically by name to ensure consistent selection
		sort.Slice(listResp.Data, func(i, j int) bool {
			return listResp.Data[i].Name < listResp.Data[j].Name
		})

		// Return the first stack after sorting
		stack = listResp.Data[0]
	}

	// Populate all fields
	data.ID = types.StringValue(stack.Id)
	data.OrganizationID = types.StringValue(stack.OrganizationId)
	data.Name = types.StringValue(stack.Name)
	data.RegionID = types.StringValue(stack.RegionID)
	data.Status = types.StringValue(stack.Status)
	data.State = types.StringValue(stack.State)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
