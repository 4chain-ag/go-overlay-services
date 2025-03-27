package commands_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/assert"
)

// This is an example describing how the handler unit tests
// should be structured and tested based on the HTTP standard package.
func TestSubmitTransactionHandler_Handle(t *testing.T) {
	// given:
	app := commands.NewSubmitTransactionCommandHandler(submitTransactionProviderAlwaysOK{})
	ts := httptest.NewServer(http.HandlerFunc(app.Handle))
	defer ts.Close()

	// when:
	res, err := ts.Client().Post(ts.URL, "application/json", nil)

	// then:
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.StatusCode, http.StatusCreated)
	defer res.Body.Close()

	var actual commands.SubmitTransactionHandlerResponse
	expected := commands.SubmitTransactionHandlerResponse{Steak: overlay.Steak{}}

	assert.NoError(t, jsonutil.DecodeResponseBody(res, &actual))
	assert.Equal(t, expected, actual)
}

type submitTransactionProviderAlwaysOK struct{}

func (submitTransactionProviderAlwaysOK) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {
	return overlay.Steak{}, nil
}
