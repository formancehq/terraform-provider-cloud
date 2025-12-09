package datasources

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/pkg/tracing"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	_ datasource.DataSource                     = &DatasourcesTracer{}
	_ datasource.DataSourceWithConfigure        = &DatasourcesTracer{}
	_ datasource.DataSourceWithConfigValidators = &DatasourcesTracer{}
)

type DatasourcesTracer struct {
	tracer          trace.Tracer
	logger          logging.Logger
	underlyingValue any
}

func injectTraceContext(ctx context.Context, res any, funcName string) context.Context {
	name := reflect.TypeOf(res).Elem().Name()
	ctx = logging.ContextWithField(ctx, "resource", strings.ToLower(name))
	ctx = logging.ContextWithField(ctx, "operation", strings.ToLower(funcName))

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ctx
	}

	// TODO: implement a logger hook to automatically add trace context to logs
	headerCarrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, headerCarrier)
	for k, v := range headerCarrier {
		ctx = logging.ContextWithField(ctx, k, v)
	}

	span.SetAttributes(
		attribute.String("resource", strings.ToLower(name)),
		attribute.String("operation", strings.ToLower(funcName)),
	)
	return ctx
}

var (
	ErrConfigure = fmt.Errorf("error in configure")
	ErrRead      = fmt.Errorf("error in read")
	ErrSchema    = fmt.Errorf("error in schema")
	ErrMetadata  = fmt.Errorf("error in metadata")
	ErrValidate  = fmt.Errorf("error in validate")
)

// Configure implements datasource.DataSourceWithConfigure.
func (d *DatasourcesTracer) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	ctx = logging.ContextWithLogger(ctx, d.logger)
	operation := "Configure"
	if v, ok := d.underlyingValue.(datasource.DataSourceWithConfigure); ok {
		_ = tracing.TraceError(ctx, d.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, d.underlyingValue, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debugf("completed")
			v.Configure(ctx, req, res)
			if res.Diagnostics.HasError() {
				return ErrConfigure
			}
			return nil
		})
	}
}

// ConfigValidators implements datasource.DataSourceWithConfigValidators.
func (d *DatasourcesTracer) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	ctx = logging.ContextWithLogger(ctx, d.logger)
	operation := "ConfigValidators"
	if v, ok := d.underlyingValue.(datasource.DataSourceWithConfigValidators); ok {
		ret, _ := tracing.Trace(ctx, d.tracer, operation, func(ctx context.Context) ([]datasource.ConfigValidator, error) {
			ctx = injectTraceContext(ctx, d.underlyingValue, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debugf("completed")
			return v.ConfigValidators(ctx), nil
		})
		return ret
	}
	return nil
}

// Metadata implements datasource.DataSourceWithConfigValidators.
func (d *DatasourcesTracer) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	ctx = logging.ContextWithLogger(ctx, d.logger)
	operation := "Metadata"
	if v, ok := d.underlyingValue.(datasource.DataSource); ok {
		_ = tracing.TraceError(ctx, d.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, d.underlyingValue, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debugf("completed")
			v.Metadata(ctx, req, res)
			return nil
		})
	}
}

// Read implements datasource.DataSourceWithConfigValidators.
func (d *DatasourcesTracer) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	ctx = logging.ContextWithLogger(ctx, d.logger)
	operation := "Read"
	if v, ok := d.underlyingValue.(datasource.DataSource); ok {
		_ = tracing.TraceError(ctx, d.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, d.underlyingValue, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debugf("completed")
			v.Read(ctx, req, res)
			if res.Diagnostics.HasError() {
				return ErrRead
			}
			return nil
		})
	}
}

// Schema implements datasource.DataSourceWithConfigValidators.
func (d *DatasourcesTracer) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	ctx = logging.ContextWithLogger(ctx, d.logger)
	operation := "Schema"
	if v, ok := d.underlyingValue.(datasource.DataSource); ok {
		_ = tracing.TraceError(ctx, d.tracer, operation, func(ctx context.Context) error {
			ctx = injectTraceContext(ctx, d.underlyingValue, operation)
			logging.FromContext(ctx).Debug("call")
			defer logging.FromContext(ctx).Debugf("completed")
			v.Schema(ctx, req, res)
			if res.Diagnostics.HasError() {
				return ErrSchema
			}
			return nil
		})
	}
}

func NewDatasourcesTracer(tracer trace.Tracer, logger logging.Logger, res func() datasource.DataSource) *DatasourcesTracer {
	return &DatasourcesTracer{
		tracer:          tracer,
		logger:          logger,
		underlyingValue: res(),
	}
}
