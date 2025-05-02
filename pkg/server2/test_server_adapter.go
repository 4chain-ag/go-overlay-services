package server2

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
)

// TestHTTPServer is a helper type for testing HTTP handlers using the Fiber app directly.
// It implements http.RoundTripper so it can be plugged into HTTP clients like resty for in-process testing
// without the need to start a real network server.
type TestHTTPServer struct {
	t       *testing.T  // The test instance used for asserting test failures.
	srv     *ServerHTTP // The Fiber server instance under test.
	timeout int         // Request timeout used by Fiber's Test method (-1 for no timeout).
}

// RoundTrip implements the http.RoundTripper interface.
// It allows http.Client-compatible tools (like resty) to execute requests against the in-memory Fiber app.
func (ts *TestHTTPServer) RoundTrip(req *http.Request) (*http.Response, error) {
	ts.t.Helper()
	return ts.srv.app.Test(req, ts.timeout)
}

// NewTestHTTPServer creates a new TestHTTPServer instance.
// It initializes the underlying Fiber server using the provided ServerOption values.
func NewTestHTTPServer(t *testing.T, opts ...ServerOption) *TestHTTPServer {
	return &TestHTTPServer{
		t:       t,
		timeout: -1,
		srv:     New(opts...),
	}
}

// TestFixture is a test utility structure that wraps a TestHTTPServer and provides
// an HTTP client (resty) configured to send requests to the in-memory Fiber app.
type TestFixture struct {
	t   *testing.T      // The test instance used for assertions and logging.
	srv *TestHTTPServer // The in-process test server.
}

// Client returns a new resty.Client instance configured to route all requests
// to the in-memory test server. Any HTTP error triggers a test failure.
func (f *TestFixture) Client() *resty.Client {
	f.t.Helper()

	c := resty.New()
	c.OnError(func(r *resty.Request, err error) {
		f.t.Fatalf("HTTP request ended with unexpected error: %v", err)
	})
	c.GetClient().Transport = f.srv

	return c
}

// NewTestFixture creates a new TestFixture with the given ServerOption values.
// It sets up an in-process test server and prepares it for testing with resty.
func NewTestFixture(t *testing.T, opts ...ServerOption) *TestFixture {
	return &TestFixture{
		t:   t,
		srv: NewTestHTTPServer(t, opts...),
	}
}
