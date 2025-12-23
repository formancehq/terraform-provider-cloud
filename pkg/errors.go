package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/sdkerrors"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.opentelemetry.io/otel/trace"
)

type Error struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func HandleSDKError(ctx context.Context, err error, diag *diag.Diagnostics) {
	sharedError := &Error{
		ErrorCode:    "INTERNAL",
		ErrorMessage: err.Error(),
	}

	tmp := &sdkerrors.SDKError{}
	if errors.As(err, &tmp) {
		err = errors.New(tmp.Body)
	}

	errResponse := &Error{}
	if e := json.Unmarshal([]byte(err.Error()), errResponse); e == nil {
		sharedError = errResponse
	}
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		traceparent := fmt.Sprintf("%s-%s", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		sharedError.ErrorMessage = fmt.Sprintf("[Traceparent: %s] %s", traceparent, sharedError.ErrorMessage)
	}
	diag.AddError(
		string(sharedError.ErrorCode),
		sharedError.ErrorMessage,
	)
}
