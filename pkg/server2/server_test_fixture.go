package server2

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

// fiberRoundTripper is a custom http.RoundTripper that routes requests through the in-memory Fiber test server.
// It is used to simulate HTTP requests directly against the server without starting a real network listener.
type fiberRoundTripper struct {
	t       *testing.T
	srv     *ServerHTTP
	timeout int
}

// RoundTrip executes an HTTP request using Fiber’s in-memory test engine.
// It returns the response or fails the test if an error occurs.
func (f *fiberRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	f.t.Helper()
	return f.srv.app.Test(req, f.timeout)
}

// ServerTestFixture provides an isolated test fixture for executing HTTP requests
// against a fully initialized server instance using in-memory transport.
type ServerTestFixture struct {
	t            *testing.T
	roundTripper http.RoundTripper
}

// Client returns a preconfigured Resty client that uses the in-memory test server.
// It automatically fails the test on unexpected transport errors.
func (f *ServerTestFixture) Client() *resty.Client {
	f.t.Helper()

	c := resty.New()
	c.OnError(func(r *resty.Request, err error) {
		require.NoError(f.t, err, "HTTP request ended with unexpected error")
	})
	c.GetClient().Transport = f.roundTripper

	return c
}

// NewServerTestFixture creates a new test fixture with a fully initialized server instance
// and a custom in-memory HTTP round tripper. Panics if server initialization fails.
func NewServerTestFixture(t *testing.T, opts ...ServerOption) *ServerTestFixture {
	return &ServerTestFixture{
		t: t,
		roundTripper: &fiberRoundTripper{
			t:       t,
			timeout: -1,
			srv:     New(opts...),
		},
	}
}
