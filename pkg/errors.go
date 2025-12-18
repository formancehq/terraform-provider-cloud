package pkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/sdkerrors"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.opentelemetry.io/otel/trace"
)

func HandleSDKError(ctx context.Context, err error, diag *diag.Diagnostics) {

	sharedErr := &sdkerrors.Error{
		ErrorCode:    "INTERNAL",
		ErrorMessage: pointer.For("unexpected error"),
	}

	if errors.Is(err, sharedErr) {
		sharedErr = err.(*sdkerrors.Error)
	}

	msg := ""
	if sharedErr.ErrorMessage != nil {
		msg = *sharedErr.ErrorMessage
	}
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		diag.AddError("traceparent", fmt.Sprintf("%s-%s", span.SpanContext().TraceID(), span.SpanContext().SpanID()))
	}

	diag.AddError(sharedErr.ErrorCode, msg)
}
