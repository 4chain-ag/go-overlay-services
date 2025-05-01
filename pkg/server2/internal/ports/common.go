package ports

import (
	"fmt"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
)

// RequestTimeout defines the default duration after which a request is considered timed out.
const RequestTimeout = 5 * time.Second

// RequestBodyLimit1GB defines the default maximum size for request bodies (1 GB).
const RequestBodyLimit1GB = 1000 * 1024 * 1024

// NewRequestBodyTooLargeResponse creates a bad request response when the submitted request body exceeds the allowed size.
// It takes the maximum allowed size, and returns an openapi.BadRequestResponse with a message indicating that the request body is too large.
func NewRequestBodyTooLargeResponse(limit int64) openapi.BadRequestResponse {
	return openapi.BadRequestResponse{
		Details: &map[string]any{"bytes_read_limit": limit},
		Message: "The submitted octet-stream exceeds the maximum allowed size",
	}
}

// NewRequestTimeoutResponse creates a timeout response when a request exceeds the allowed time limit.
// It takes the timeout duration and returns an openapi.RequestTimeoutResponse with a message indicating that the request timed out.
func NewRequestTimeoutResponse(limit time.Duration) openapi.RequestTimeoutResponse {
	return openapi.BadRequestResponse{
		Message: fmt.Sprintf("The submitted request exceeded the timeout limit of %d seconds.", int(limit.Seconds())),
	}
}

// NewRequestMissingHeaderResponse creates a bad request response indicating that a required HTTP header is missing.
// It takes the name of the missing header and returns an openapi.BadRequestResponse with a descriptive message.
func NewRequestMissingHeaderResponse(header string) openapi.BadRequestResponse {
	return openapi.BadRequestResponse{
		Message: fmt.Sprintf("The submitted request does not include required header: %s.", header),
	}
}
