package overlayhttp

import (
	"fmt"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
)

// NewRequestBodyTooLargeResponse creates a bad request response when the submitted request body exceeds the allowed size.
// It takes the actual size of the request body and the maximum allowed size, and returns an openapi.BadRequestResponse
// with a message indicating that the request body is too large.
func NewRequestBodyTooLargeResponse(actual int64, limit int64) openapi.BadRequestResponse {
	return openapi.Error{
		Details: &map[string]any{"bytes_read": actual},
		Message: fmt.Sprintf("The submitted octet-stream exceeds the maximum allowed size of %d bytes.", limit),
	}
}

// NewRequestTimeoutResponse creates a timeout response when a request exceeds the allowed time limit.
// It takes the timeout duration and returns an openapi.RequestTimeoutResponse with a message indicating that the request timed out.
func NewRequestTimeoutResponse(limit time.Duration) openapi.RequestTimeoutResponse {
	return openapi.Error{Message: fmt.Sprintf("The submitted request exceeded the timeout limit of %d seconds.", int(limit.Seconds()))}
}
