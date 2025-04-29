package server2

import "net/http"

// ServerTestAdapter is a test utility that wraps a ServerHTTP instance,
// allowing HTTP requests to be tested using the Fiber app's internal Test method.
type ServerTestAdapter struct {
	srv *ServerHTTP
}

// TestRequest performs an HTTP request against the underlying Fiber app instance.
// It accepts an *http.Request and a timeout (in milliseconds), and returns the *http.Response or an error.
func (t *ServerTestAdapter) TestRequest(req *http.Request, timeout int) (*http.Response, error) {
	return t.srv.app.Test(req, timeout)
}

// NewServerTestAdapter creates a new instance of ServerTestAdapter with the given ServerOption(s).
// It initializes the ServerHTTP instance using the same configuration logic as production,
// making it suitable for use in integration or functional tests.
func NewServerTestAdapter(opts ...ServerOption) *ServerTestAdapter {
	return &ServerTestAdapter{srv: New(opts...)}
}
