package testabilities

import (
	"encoding/json"
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
