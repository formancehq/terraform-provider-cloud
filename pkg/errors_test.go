package pkg

import (
	"context"
	"errors"
	"testing"

	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/sdkerrors"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/require"
)

func TestHandleSDKError(t *testing.T) {
	for _, tt := range []struct {
		name     string
		err      error
		expected diag.Diagnostic
	}{
		{
			name: "Test case 1",
			err:  errors.New(""),
			expected: diag.NewErrorDiagnostic(
				"INTERNAL",
				"",
			),
		},
		{
			name: "Error string",
			err:  errors.New(`{"errorCode":"VALIDATION","errorMessage":"invalid config: polling period invalid: polling period cannot be lower than minimum of 20m0s: validation error: validation error"}`),
			expected: diag.NewErrorDiagnostic(
				"VALIDATION",
				"invalid config: polling period invalid: polling period cannot be lower than minimum of 20m0s: validation error: validation error",
			),
		},
		{
			name: "SDKError case",
			err: &sdkerrors.SDKError{
				Body: `{"errorCode":"SOME_ERROR","errorMessage":"An error occurred"}`,
			},
			expected: diag.NewErrorDiagnostic("SOME_ERROR", "An error occurred"),
		},
		{
			name:     "invalid error type",
			err:      errors.New("some random error"),
			expected: diag.NewErrorDiagnostic("INTERNAL", "some random error"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			diag := make(diag.Diagnostics, 0)
			HandleSDKError(context.Background(), tt.err, &diag)
			require.Len(t, diag, 1)
			require.Equal(t, tt.expected, diag[0])
		})
	}
}
