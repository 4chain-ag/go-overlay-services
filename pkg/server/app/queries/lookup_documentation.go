package queries

import (
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/gofiber/fiber/v2"
)

// LookupDocumentationProvider defines the contract that must be fulfilled
// to send a lookup service documentation request to the overlay engine for further processing.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type LookupDocumentationProvider interface {
	GetDocumentationForLookupServiceProvider(lookupService string) (string, error)
}

// LookupDocumentationHandler orchestrates the processing flow of a lookup documentation
// request, including the request parameter validation, converting the request
// into an overlay-engine-compatible format, and applying any other necessary
// logic before invoking the engine. It returns the requested lookup service
// documentation in the text/markdown format.
type LookupDocumentationHandler struct {
	provider LookupDocumentationProvider
}

// Handle orchestrates the processing flow of a lookup documentation request.
// It extracts the lookupService query parameter, invokes the engine provider,
// and returns the documentation as markdown with the appropriate status code.
func (l *LookupDocumentationHandler) Handle(c *fiber.Ctx) error {
	lookupService := c.Query("lookupService")	
	documentation, err := l.provider.GetDocumentationForLookupServiceProvider(lookupService)
	if err != nil {
		if inner := c.Status(fiber.StatusBadRequest).JSON(dto.HandlerResponseNonOK); inner != nil {
			return fmt.Errorf("failed to send JSON response: %w", inner)
		}
		return nil
	}

	// Set Content-Type header to text/markdown
	c.Set("Content-Type", "text/markdown")
	if err := c.Status(fiber.StatusOK).Send([]byte(documentation)); err != nil {
		return fmt.Errorf("failed to send markdown response: %w", err)
	}
	
	return nil
}

// NewLookupDocumentationHandler returns an instance of a LookupDocumentationHandler, utilizing
// an implementation of LookupDocumentationProvider. If the provided argument is nil, it triggers a panic.
func NewLookupDocumentationHandler(provider LookupDocumentationProvider) *LookupDocumentationHandler {
	if provider == nil {
		panic("lookup documentation provider is nil")
	}
	return &LookupDocumentationHandler{
		provider: provider,
	}
}
