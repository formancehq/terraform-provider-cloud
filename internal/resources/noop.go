package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ resource.Resource              = &Noop{}
	_ resource.ResourceWithConfigure = &Noop{}
)

type Noop struct {
	logger logging.Logger
	store  *internal.Store
}

type NoopModel struct{}

func NewNoop(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &Noop{
			logger: logger,
		}
	}
}

var SchemaNoop = schema.Schema{
	Description: "",
	Attributes:  map[string]schema.Attribute{},
}

// Schema implements resource.Resource.
func (s *Noop) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaNoop
}

// Configure implements resource.ResourceWithConfigure.
func (s *Noop) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	store, ok := req.ProviderData.(*internal.Store)
	if !ok {
		res.Diagnostics.AddError(
			ErrProviderDataNotSet.Error(),
			fmt.Sprintf("Expected *internal.Store, got: %T", req.ProviderData),
		)
		return
	}

	s.store = store
}

// Create implements resource.Resource.
func (s *Noop) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan NoopModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *Noop) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state NoopModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

}

// Metadata implements resource.Resource.
func (s *Noop) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_noop"
}

// Read implements resource.Resource.
func (s *Noop) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state NoopModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *Noop) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var state NoopModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}
