package jsonutil

import (
	"bytes"
	"encoding/json"
	"errors"
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

// DecodeRequestBody attempts to decode the HTTP request body into given destination
// argument. It returns an error if the internal decoding operation fails; otherwise,
// it returns nil, indicating successful processing.
func DecodeRequestBody(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dst)
	if err != nil {
		return fmt.Errorf("decoding http request body op failure: %w", err)
	}
	return nil
}

// DecodeBytes deserializes JSON bytes into the provided destination object.
// It uses a json.Decoder to parse the byte array and returns an error
// joined with JSONDecoderFailure if the decoding process fails.
func DecodeBytes(bb []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewBuffer(bb))
	err := dec.Decode(dst)
	if err != nil {
		return errors.Join(err, JSONDecoderFailure)
	}
	return nil
}

// JSONDecoderFailure represents an error that occurs when the JSON decoder fails
// to parse the input data into the expected structure.
var JSONDecoderFailure = errors.New("failed to decode JSON payload")
