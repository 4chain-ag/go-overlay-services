package commands_test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/commands"
	"github.com/stretchr/testify/require"
)

func TestARCIngestHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		arcAPIKey      string
		expectedStatus int
	}{
		{
			name:           "should success with 200 when ARC API key is configured",
			arcAPIKey:      "3988c27c-b9ab-4f3e-8b12-d3dc7218eb2d",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "should fail with 500 when no ARC API key is configured",
			arcAPIKey:      "",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	rand.Shuffle(len(tests), func(i, j int) {
		tests[i], tests[j] = tests[j], tests[i]
	})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given:
			handler, err := commands.NewARCIngestHandler(tc.arcAPIKey)
			require.NoError(t, err)

			ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
			defer ts.Close()

			// when:
			res, err := ts.Client().Post(ts.URL, "application/json", nil)

			// then:
			require.NoError(t, err)
			require.NotNil(t, res)

			require.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}
