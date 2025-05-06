package ports

import (
	"fmt"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
)

// RequestTimeout defines the default duration after which a request is considered timed out.
const RequestTimeout = 5 * time.Second

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

// UnhandledErrorTypeResponse is the default response returned when an error occurs
// that does not match any known or handled ErrorType.
// It represents a generic internal server error to avoid exposing internal details to the client.
var UnhandledErrorTypeResponse = openapi.InternalServerErrorResponse{
	Message: "An internal error occurred during processing the request. Please try again later or contact the support team.",
}
