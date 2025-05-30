package ports

import (
	"context"
	"encoding/json"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/gofiber/fiber/v2"
)

// LookupQuestionService defines the contract for handling lookup questions.
type LookupQuestionService interface {
	LookupQuestion(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
}

// LookupQuestionHandler handles lookup question requests.
type LookupQuestionHandler struct {
	service LookupQuestionService
}

// Handle processes a lookup question request and returns the answer.
func (h *LookupQuestionHandler) Handle(c *fiber.Ctx, params openapi.LookupQuestionBody) error {
	if err := c.BodyParser(&params); err != nil {
		return app.NewLookupQuestionInvalidRequestBodyResponse()
	}

	// Convert the Query map to JSON for the lookup question
	var queryJSON json.RawMessage
	if params.Query != nil {
		if queryBytes, err := json.Marshal(params.Query); err != nil {
			return app.NewLookupQuestionInvalidRequestBodyResponse()
		} else {
			queryJSON = queryBytes
		}
	}

	question := lookup.LookupQuestion{
		Service: params.Service,
		Query:   queryJSON,
	}

	answer, err := h.service.LookupQuestion(c.UserContext(), &question)
	if err != nil {
		return err
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

	return &LookupQuestionHandler{service: app.NewLookupQuestionService(provider)}
}

// NewLookupQuestionSuccessResponse creates a successful response for a lookup question.
func NewLookupQuestionSuccessResponse(answer *lookup.LookupAnswer) *openapi.LookupAnswer {
	var outputs []openapi.OutputListItem
	if len(answer.Outputs) > 0 {
		outputs = make([]openapi.OutputListItem, len(answer.Outputs))
		for i, output := range answer.Outputs {
			outputs[i] = openapi.OutputListItem{
				Beef:        output.Beef,
				OutputIndex: output.OutputIndex,
			}
		}
	}

	var resultStr string
	if answer.Result != nil {
		if resultBytes, err := json.Marshal(answer.Result); err == nil {
			resultStr = string(resultBytes)
		}
	}

	return &openapi.LookupAnswer{
		Outputs: outputs,
		Result:  resultStr,
		Type:    string(answer.Type),
	}
}
