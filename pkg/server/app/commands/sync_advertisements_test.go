package commands_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/stretchr/testify/assert"
)

// This is an example describing how the handler unit tests
// should be structured and tested based on the HTTP standard package.
func TestSyncAdvertisementsHandler_Handle(t *testing.T) {
	// given:
	app := commands.NewSyncAdvertisementsCommandHandler(syncAdvertisementsProviderAlwaysOK{})
	ts := httptest.NewServer(http.HandlerFunc(app.Handle))
	defer ts.Close()

	// when:
	res, err := ts.Client().Post(ts.URL, "application/json", nil)

	// then:
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()

	var actual commands.SyncAdvertisementsHandlerResponse
	expected := commands.SyncAdvertisementsHandlerResponse{Message: "OK"}

	assert.NoError(t, jsonutil.DecodeResponseBody(res, &actual))
	assert.Equal(t, expected, actual)
}

// Stubbed provider that always succeeds
type syncAdvertisementsProviderAlwaysOK struct{}

func (syncAdvertisementsProviderAlwaysOK) SyncAdvertisements() error {
	return nil
}
