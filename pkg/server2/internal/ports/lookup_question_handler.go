package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// LookupQuestionService defines the interface for a service that performs
// lookup operations based on a service name and query parameters. It returns
// the result in a LookupAnswerDTO format suitable for further processing or response.
type LookupQuestionService interface {
	LookupQuestion(ctx context.Context, service string, query map[string]any) (*app.LookupAnswerDTO, error)
}

// LookupQuestionHandler is an HTTP handler that processes requests to perform
// a lookup operation on a specific question. It uses a LookupQuestionService to
// evaluate the query and returns the results formatted according to the OpenAPI schema.
type LookupQuestionHandler struct {
	service LookupQuestionService
}

// Handle is the HTTP endpoint function for handling a lookup question request.
// It parses the request body, validates and forwards the data to the service layer,
// and returns a JSON response with the lookup results or an appropriate error.
func (h *LookupQuestionHandler) Handle(c *fiber.Ctx, params openapi.LookupQuestionBody) error {
	err := c.BodyParser(&params)
	if err != nil {
		return NewRequestBodyParserError(err)
	}

	dto, err := h.service.LookupQuestion(c.UserContext(), params.Service, params.Query)
	if err != nil {
		return err
	}

	res, err := NewLookupQuestionSuccessResponse(dto)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

// NewLookupQuestionHandler creates and returns a new LookupQuestionHandler.
// It initializes the handler with a LookupQuestionService using the provided provider.
// The function panics if the provider is nil.
func NewLookupQuestionHandler(provider app.LookupQuestionProvider) *LookupQuestionHandler {
	if provider == nil {
		panic("lookup question provider is nil")
	}
	return &LookupQuestionHandler{service: app.NewLookupQuestionService(provider)}
}

// NewLookupQuestionSuccessResponse constructs an OpenAPI-compatible LookupAnswer
// from a LookupAnswerDTO. It marshals the output items and the result string into
// the format expected by the HTTP client.
func NewLookupQuestionSuccessResponse(dto *app.LookupAnswerDTO) (*openapi.LookupAnswer, error) {
	var outputs []openapi.OutputListItem
	if len(dto.Outputs) > 0 {
		outputs = make([]openapi.OutputListItem, len(dto.Outputs))
		for i, output := range dto.Outputs {
			outputs[i] = openapi.OutputListItem{
				Beef:        output.BEEF,
				OutputIndex: output.OutputIndex,
			}
		}
	}

	return &openapi.LookupAnswer{
		Outputs: outputs,
		Result:  dto.Result,
		Type:    dto.Type,
	}, nil
}

// NewRequestBodyParserError defines the error that should be returned when the request body content
// cannot be parsed successfully due to an invalid format or any other issue that occurred during decoding phase.
func NewRequestBodyParserError(err error) app.Error {
	return app.NewRawDataProcessingError(err.Error(), "Unable to process the request body content due to an internal error. Please verify content, try again later or contact the support team.")
}
