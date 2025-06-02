package app

import (
	"context"
	"encoding/json"

	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
)

// OutputListItemDTO represents an output item in the final answer returned to the caller.
// It includes the raw output data ('BEEF') and its positional index.
type OutputListItemDTO struct {
	BEEF        []byte
	OutputIndex uint32
}

// LookupAnswerDTO is a response DTO that represents the result of a successful lookup question
// evaluation. It contains the result, output items, and the answer type.
type LookupAnswerDTO struct {
	Outputs []OutputListItemDTO
	Result  string
	Type    string
}

// LookupQuestionProvider defines the interface for a provider capable of processing lookup questions.
// It encapsulates the logic required to evaluate a question and produce a corresponding answer.
type LookupQuestionProvider interface {
	Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
}

// LookupQuestionService provides a higher-level abstraction over a LookupQuestionProvider.
// It performs validation on incoming questions before delegating the lookup operation to
// the underlying provider.
type LookupQuestionService struct {
	provider LookupQuestionProvider
}

// LookupQuestion validates the service and query fields, serializes the query to JSON,
// constructs a lookup question, delegates the call to the underlying provider, and
// converts the provider's response into a DTO.
func (s *LookupQuestionService) LookupQuestion(ctx context.Context, service string, query map[string]any) (*LookupAnswerDTO, error) {
	if len(service) == 0 {
		return nil, NewIncorrectInputWithFieldError("service")
	}
	if len(query) == 0 {
		return nil, NewIncorrectInputWithFieldError("query")
	}
	bb, err := json.Marshal(query)
	if err != nil {
		return nil, NewLookupQuestionQueryParserError(err)
	}

	answer, err := s.provider.Lookup(ctx, &lookup.LookupQuestion{
		Service: service,
		Query:   json.RawMessage(bb),
	})
	if err != nil {
		return nil, NewLookupQuestionProviderError(err)
	}

	return NewLookupAnswerDTO(answer)
}

// NewLookupQuestionService constructs a new LookupQuestionService using the provided provider.
// It panics if the given provider is nil, ensuring proper service initialization.
func NewLookupQuestionService(provider LookupQuestionProvider) *LookupQuestionService {
	if provider == nil {
		panic("lookup question provider is nil")
	}
	return &LookupQuestionService{provider: provider}
}

// NewLookupAnswerDTO transforms a LookupAnswer into a LookupAnswerDTO.
func NewLookupAnswerDTO(answer *lookup.LookupAnswer) (*LookupAnswerDTO, error) {
	var outputs []OutputListItemDTO
	if len(answer.Outputs) > 0 {
		outputs = make([]OutputListItemDTO, len(answer.Outputs))
		for i, output := range answer.Outputs {
			outputs[i] = OutputListItemDTO{
				BEEF:        output.Beef,
				OutputIndex: output.OutputIndex,
			}
		}
	}

	var result string
	if answer.Result != nil {
		bb, err := json.Marshal(answer.Result)
		if err != nil {
			return nil, NewRawDataProcessingError(err.Error(), "Unable to create the lookup question response due to an internal error. Please try again later or contact the support team.")
		}
		result = string(bb)
	}

	return &LookupAnswerDTO{
		Outputs: outputs,
		Result:  result,
		Type:    string(answer.Type),
	}, nil
}

// NewLookupQuestionQueryParserError defines the error returned when the query parameters cannot
// be parsed due to an invalid format or an internal JSON encoding error.
func NewLookupQuestionQueryParserError(err error) Error {
	return NewRawDataProcessingError(err.Error(), "Unable to process the request query params content due to an internal error. Please verify the content, try again later, or contact the support team.")
}

// NewLookupQuestionProviderError wraps a provider-level error that occurred during the lookup question processing.
// It returns a generic provider failure error message intended for client-facing responses while preserving
// the original error message internally.
func NewLookupQuestionProviderError(err error) Error {
	return NewProviderFailureError(err.Error(),
		"Unable to process lookup question due to an internal error. Please try again later or contact the support team.",
	)
}
