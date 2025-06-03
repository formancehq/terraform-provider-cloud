package pkg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func HandleSDKError(ctx context.Context, res *http.Response, diag *diag.Diagnostics) {
	var details []string
	traceId := res.Header.Get("X-Trace-Id")
	if traceId != "" {
		details = append(details, fmt.Sprintf("Trace ID: %s", traceId))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		details = append(details, fmt.Sprintf("Error reading response body: %s", err.Error()))
		diag.AddError("Unexpected error", strings.Join(details, "\r\n"))
		return
	}

	errResp := sdk.Error{}
	if err := errResp.UnmarshalJSON(body); err != nil {
		details = append(details, fmt.Sprintf("Error unmarshalling error response: %s", err.Error()))
		details = append(details, fmt.Sprintf("Body: %s", body))
		diag.AddError("Unexpected error", strings.Join(details, "\r\n"))
		return
	}

	if errResp.ErrorMessage != nil {
		details = append(details, fmt.Sprintf("Message: %s", *errResp.ErrorMessage))
	}
	diag.AddError(fmt.Sprintf("Error %s", errResp.ErrorCode), strings.Join(details, "\r\n"))
}
