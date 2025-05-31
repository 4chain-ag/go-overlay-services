package app

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
)

// LookupQuestionProvider defines the contract that must be fulfilled to process lookup questions in the overlay engine.
type LookupQuestionProvider interface {
	Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
}

// LookupQuestionService provides functionality for processing lookup questions.
type LookupQuestionService struct {
	provider LookupQuestionProvider
}

// LookupQuestion processes a lookup question and returns the answer.
func (s *LookupQuestionService) LookupQuestion(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	if question == nil {
		return nil, NewInvalidLookupQuestionError()
	}

	if question.Service == "" {
		return nil, NewLookupQuestionMissingServiceFieldError()
	}

	if question.Query == nil {
		return nil, NewLookupQuestionMissingQueryFieldError()
	}

	answer, err := s.provider.Lookup(ctx, question)
	if err != nil {
		return nil, NewLookupQuestionProviderError(err)
	}

	return answer, nil
}

// NewLookupQuestionService creates a new LookupQuestionService instance using the given provider.
// Panics if the provider is nil.
func NewLookupQuestionService(provider LookupQuestionProvider) *LookupQuestionService {
	if provider == nil {
		panic("lookup question provider is nil")
	}

	return &LookupQuestionService{
		provider: provider,
	}
}

// NewInvalidLookupQuestionError returns an Error indicating that the lookup question is nil.
func NewInvalidLookupQuestionError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "lookup question cannot be nil",
		slug:      "The lookup question must be provided and cannot be nil.",
	}
}

// NewLookupQuestionMissingServiceFieldError returns an Error indicating that the service field is missing.
func NewLookupQuestionMissingServiceFieldError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "missing required service field in the request",
		slug:      "The service field is required in the lookup question request.",
	}
}

// NewLookupQuestionProviderError returns an Error indicating that the overlay engine failed to process a lookup question.
func NewLookupQuestionProviderError(err error) Error {
	return Error{
		errorType: ErrorTypeProviderFailure,
		err:       err.Error(),
		slug:      "Unable to process lookup question due to an error in the overlay engine.",
	}
}

// NewLookupQuestionInvalidRequestBodyResponse creates an error response for invalid request body.
func NewLookupQuestionInvalidRequestBodyResponse() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "invalid request body format or structure",
		slug:      "The request body must be a valid JSON object with a 'service' field and a 'query' field.",
	}
}

// NewLookupQuestionMissingQueryFieldError returns an Error indicating that the query field is invalid.
func NewLookupQuestionMissingQueryFieldError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "invalid query field format or structure",
		slug:      "The query field must be a valid JSON object.",
	}
}
