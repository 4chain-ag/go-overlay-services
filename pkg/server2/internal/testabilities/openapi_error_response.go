package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
)

// NewTestOpenapiErrorResponse creates an openapi.Error response from the given app.Error,
// primarily for use in tests. It sets the error message to the error's slug.
func NewTestOpenapiErrorResponse(t *testing.T, err app.Error) openapi.Error {
	t.Helper()
	return openapi.Error{
		Message: err.Slug(),
	}
}
