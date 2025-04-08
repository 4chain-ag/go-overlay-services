package testutil

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// SubmitTransactionProviderAlwaysSuccess mocks a transaction provider that always succeeds

type SubmitTransactionProviderAlwaysSuccess struct{ ExpectedSteak overlay.Steak }

// NewSubmitTransactionProviderAlwaysSuccess creates a new instance of SubmitTransactionProviderAlwaysSuccess

func NewSubmitTransactionProviderAlwaysSuccess(steak overlay.Steak) *SubmitTransactionProviderAlwaysSuccess {

	return &SubmitTransactionProviderAlwaysSuccess{ExpectedSteak: steak}

}

// Submit implements the transaction submission interface

func (s SubmitTransactionProviderAlwaysSuccess) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {

	// Call the onSteakReady callback to simulate async completion

	onSteakReady(&s.ExpectedSteak)

	return nil, nil

}

// SubmitTransactionProviderAlwaysFailure mocks a transaction provider that always fails

type SubmitTransactionProviderAlwaysFailure struct{ ExpectedErr error }

// Submit implements the transaction submission interface

func (s SubmitTransactionProviderAlwaysFailure) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {

	return nil, s.ExpectedErr

}

// NewSubmitTransactionProviderAlwaysFailure creates a new instance of SubmitTransactionProviderAlwaysFailure

func NewSubmitTransactionProviderAlwaysFailure(err error) *SubmitTransactionProviderAlwaysFailure {

	return &SubmitTransactionProviderAlwaysFailure{ExpectedErr: err}

}

// SubmitTransactionProviderNeverCallback mocks a transaction provider that never calls the callback

type SubmitTransactionProviderNeverCallback struct{}

// Submit implements the transaction submission interface

func (s SubmitTransactionProviderNeverCallback) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {

	// Never call the callback which then should trigger the timeout

	return overlay.Steak{}, nil

}
