package server_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	config "github.com/4chain-ag/go-overlay-services/pkg/appconfig"
	"github.com/4chain-ag/go-overlay-services/pkg/server"
)

func Test_AuthorizationBearerTokenMiddleware(t *testing.T) {
	// Given
	adminToken := "valid_admin_token"

	handler := server.AdminAuth(adminToken)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	ts := httptest.NewServer(handler)
	defer ts.Close()

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Route access without a token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Route access with an invalid token",
			token:          "invalid_token",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Route access with a valid token",
			token:          "valid_admin_token",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			resp, err := ts.Client().Do(req)

			// Then
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func Test_HTTPServer_ShouldShutdownAfterSendingInterruptSig(t *testing.T) {
	// given:
	cfg := config.Defaults()
	opts := []server.HTTPOption{
		server.WithConfig(&cfg),
	}
	httpAPI := server.New(opts...)

	// when:
	done := httpAPI.StartWithGracefulShutdown()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		slog.Info("Sending os.Interrupt signal to the HTTP API", slog.Attr{
			Key:   "process_id",
			Value: slog.IntValue(os.Getpid()),
		})

		process, err := os.FindProcess(os.Getpid())
		require.NoError(t, err, "Failed to find HTTP API process")

		require.NoError(t, process.Signal(os.Interrupt), "Failed to send os.Interrupt signal to the HTTP API")
	}()

	wg.Wait()

	// then:
	_, ok := <-done
	require.False(t, ok, "Server did not shut down as expected")
}
