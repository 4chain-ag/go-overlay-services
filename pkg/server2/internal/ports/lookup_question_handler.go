package ports

import (
	"context"
	"encoding/json"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// LookupQuestionService defines the contract for handling lookup questions.
type LookupQuestionService interface {
	LookupQuestion(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
}

// LookupQuestionHandlerRequest represents the request body for a lookup question.
type LookupQuestionHandlerRequest struct {
	Service string          `json:"service"`
	Query   json.RawMessage `json:"query"`
}

// ToLookupQuestion converts the LookupQuestionHandlerRequest to a lookup.LookupQuestion.
func (r LookupQuestionHandlerRequest) ToLookupQuestion() *lookup.LookupQuestion {
	return &lookup.LookupQuestion{
		Service: r.Service,
		Query:   r.Query,
	}
}

// LookupQuestionHandler handles lookup question requests.
type LookupQuestionHandler struct {
	service LookupQuestionService
}

// Handle processes a lookup question request and returns the answer.
func (h *LookupQuestionHandler) Handle(c *fiber.Ctx) error {
	var request LookupQuestionHandlerRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NewInvalidRequestBodyResponse())
	}

	question := request.ToLookupQuestion()
	answer, err := h.service.LookupQuestion(c.UserContext(), question)
	if err != nil {
		var appErr app.Error
		if errors.As(err, &appErr) {
			return c.Status(fiber.StatusBadRequest).JSON(openapi.Error{
				Message: appErr.Slug(),
			})
		}
		appErr = app.NewLookupQuestionProviderError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{
			Message: appErr.Slug(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(NewLookupQuestionSuccessResponse(answer))
}

// NewLookupQuestionHandler creates a new LookupQuestionHandler.
// It initializes a LookupQuestionService with the provided provider.
// Panics if the provider is nil.
func NewLookupQuestionHandler(provider app.LookupQuestionProvider) *LookupQuestionHandler {
	if provider == nil {
		panic("lookup question provider is nil")
	}

	return &LookupQuestionHandler{
		service: app.NewLookupQuestionService(provider),
	}
}

// NewLookupQuestionSuccessResponse creates a successful response for a lookup question.
func NewLookupQuestionSuccessResponse(answer *lookup.LookupAnswer) *openapi.LookupAnswer {
	answerMap := answer.Result.(map[string]interface{})
	return &openapi.LookupAnswer{
		Answer: answerMap,
	}
}

// NewInvalidRequestBodyResponse creates an error response for invalid request body.
func NewInvalidRequestBodyResponse() openapi.BadRequestResponse {
	return openapi.Error{
		Message: "Invalid request body format or structure. Please check the API documentation for the correct format.",
	}
}
