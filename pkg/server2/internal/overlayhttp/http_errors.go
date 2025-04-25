package overlayhttp

import (
	"fmt"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
)

func NewRequestBodyTooLargeResponse(actual int64, limit int64) openapi.BadRequestResponse {
	return openapi.Error{
		Details: &map[string]any{"bytes_read": actual},
		Message: fmt.Sprintf("The submitted octet-stream exceeds the maximum allowed size of %d bytes.", limit),
	}
}

func NewRequestTimeoutResponse(limit time.Duration) openapi.RequestTimeoutResponse {
	return openapi.Error{Message: fmt.Sprintf("The submitted request exceeded the timeout limit of %d seconds.", int(limit.Seconds()))}
}
