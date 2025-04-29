package testabilities

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// DecodeResponseBody attempts to decode the HTTP response body into given destination
// argument. It returns an error if the internal decoding operation fails; otherwise,
// it returns nil, indicating successful processing.
func DecodeResponseBody(t *testing.T, res *http.Response, dst any) {
	t.Helper()

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(dst)
	require.NoError(t, err, "decoding http response body op failure")
}

// RequestBody serializes the provided value into a JSON-encoded byte slice and returns it as an io.Reader.
// This is typically used in tests to create a request body for HTTP requests.
// The function ensures that marshaling succeeds; otherwise, it stops the test execution with an error.
func RequestBody(t *testing.T, v any) io.Reader {
	t.Helper()

	bb, err := json.Marshal(v)
	require.NoError(t, err, "failed to marshal request body")
	return bytes.NewReader(bb)
}
