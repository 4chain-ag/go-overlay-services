package server_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/config"
	"github.com/stretchr/testify/require"
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

func Test_ARCAuthMiddleware(t *testing.T) {
	// Given
	arcTests := []struct {
		name           string
		arcKey         string
		expectedStatus int
	}{
		{
			name:           "Route access with empty ARC key in config",
			arcKey:         "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Route access with valid ARC key in config",
			arcKey:         "valid_arc_key",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range arcTests {
		t.Run(tt.name, func(t *testing.T) {
			handler := server.ARCAuth(tt.arcKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			ts := httptest.NewServer(handler)
			defer ts.Close()

			// When
			req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
			require.NoError(t, err)
			resp, err := ts.Client().Do(req)

			// Then
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func Test_HTTPServer_ShouldShutdownAfterSendingInterruptSig(t *testing.T) {
	// given:
	cfg := config.NewDefault()
	opts := []server.HTTPOption{
		server.WithConfig(&cfg.Server),
	}
	httpAPI, err := server.New(opts...)
	require.NoError(t, err, "Failed to create HTTP API server")

	// when:
	done := httpAPI.StartWithGracefulShutdown(context.Background())

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

func Test_HTTPServer_ShouldShutdownAfterContextCancel(t *testing.T) {
	// given:
	cfg := config.NewDefault()
	opts := []server.HTTPOption{
		server.WithConfig(&cfg.Server),
	}

	httpAPI, err := server.New(opts...)
	require.NoError(t, err)

	// when:
	trigger := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	done := httpAPI.StartWithGracefulShutdown(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-trigger
		t.Log("Triggering context cancel")
		cancel()
	}()

	close(trigger)
	wg.Wait()
	// then:
	_, ok := <-done
	require.False(t, ok, "Server did not shut down after context cancellation")
}

func Test_HTTPServer_ShouldShutdownAfterContextTimeout(t *testing.T) {
	// Given:
	cfg := config.NewDefault()
	opts := []server.HTTPOption{
		server.WithConfig(&cfg.Server),
	}

	httpAPI, err := server.New(opts...)
	require.NoError(t, err)

	// When:
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	done := httpAPI.StartWithGracefulShutdown(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Then:
		_, ok := <-done
		require.False(t, ok, "Server did not shut down after context timeout")
	}()

	wg.Wait()
}
