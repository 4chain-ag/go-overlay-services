package jsonutil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SendHTTPResponse replies to the request with the specified HTTP status code and response body content.
// Appends the "Content-Type: application/json" header to the returned response.
// In case of an internal encoding failure, it replies to the request with a StatusInternalServerError
// message and HTTP status code. The error message is returned in plain text format.
func SendHTTPResponse(w http.ResponseWriter, code int, responseBody any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	if err := enc.Encode(responseBody); err != nil {
		SendHTTPInternalServerErrorTextResponse(w)
	}
}

// SendHTTPErrorResponse replies to the request with the specified HTTP status code and a standardized
// error message in JSON format. It uses the error format {"error": "message"}.
// In case of an internal encoding failure, it replies to the request with a StatusInternalServerError
// message and HTTP status code in plain text format.
func SendHTTPErrorResponse(w http.ResponseWriter, code int, errorMessage string) {
	SendHTTPResponse(w, code, map[string]string{
		"error": errorMessage,
	})
}

// SendHTTPInternalServerErrorTextResponse replies to the request with a StatusInternalServerError
// message and HTTP status code. It does not terminate the request; the caller must ensure no further
// writes are made to w. The error message is returned in plain text format.
func SendHTTPInternalServerErrorTextResponse(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// DecodeResponseBody attempts to decode the HTTP response body into given destination
// argument. It returns an error if the internal decoding operation fails; otherwise,
// it returns nil, indicating successful processing.
func DecodeResponseBody(res *http.Response, dst any) error {
	dec := json.NewDecoder(res.Body)
	err := dec.Decode(dst)
	if err != nil {
		return fmt.Errorf("decoding http response body op failure: %w", err)
	}
	return nil
}
