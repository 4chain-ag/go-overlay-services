package ports

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/gofiber/fiber/v2"
)

// LookupQuestionService defines the contract for handling lookup questions.

type LookupQuestionService interface {
	LookupQuestion(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
}

// LookupQuestionHandlerRequest represents the request body for a lookup question.

type LookupQuestionHandlerRequest struct {
	Service string `json:"service"`

	Query json.RawMessage `json:"query"`
}

// ToLookupQuestion converts the LookupQuestionHandlerRequest to a lookup.LookupQuestion.

func (r LookupQuestionHandlerRequest) ToLookupQuestion() *lookup.LookupQuestion {

	return &lookup.LookupQuestion{

		Service: r.Service,

		Query: r.Query,
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

	switch {

	case errors.Is(err, app.ErrMissingServiceField):

		return c.Status(fiber.StatusBadRequest).JSON(NewMissingServiceFieldResponse())

	case errors.Is(err, app.ErrInvalidLookupQuestion):

		return c.Status(fiber.StatusBadRequest).JSON(NewInvalidRequestBodyResponse())

	case err != nil:

		return c.Status(fiber.StatusInternalServerError).JSON(NewLookupQuestionProviderErrorResponse())

	default:

		return c.Status(fiber.StatusOK).JSON(NewLookupQuestionSuccessResponse(answer))

	}

}

// NewLookupQuestionHandler creates a new LookupQuestionHandler.

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

		Answer: &answerMap,
	}

}

// NewLookupQuestionProviderErrorResponse creates an error response for lookup provider errors.

func NewLookupQuestionProviderErrorResponse() openapi.InternalServerErrorResponse {

	return openapi.Error{

		Message: "Unable to process lookup question due to an error in the overlay engine.",
	}

}

// NewMissingServiceFieldResponse creates an error response for missing service field.

func NewMissingServiceFieldResponse() openapi.BadRequestResponse {

	return openapi.Error{

		Message: "Missing required service field in the request body.",
	}

}

// NewInvalidRequestBodyResponse creates an error response for invalid request body.

func NewInvalidRequestBodyResponse() openapi.BadRequestResponse {

	return openapi.Error{

		Message: "Invalid request body.",
	}

}
