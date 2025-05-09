package app

import (
	"context"
	"errors"

	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
)

var (
	// ErrInvalidLookupQuestion is returned when the lookup question is nil.
	ErrInvalidLookupQuestion = errors.New("lookup question cannot be nil")

	// ErrMissingServiceField is returned when the service field is missing in the lookup question.
	ErrMissingServiceField = errors.New("missing required service field in the request")
)

// LookupQuestionProvider defines the contract that must be fulfilled
// to process lookup questions in the overlay engine.
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
		return nil, ErrInvalidLookupQuestion
	}

	if question.Service == "" {
		return nil, ErrMissingServiceField
	}

	return s.provider.Lookup(ctx, question)
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
