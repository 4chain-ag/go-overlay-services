package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// ReadBodyLimit1GB defines the maximum allowed bytes read size (in bytes).
// This limit is set to 1GB to protect against excessively large payloads.
const ReadBodyLimit1GB = 1000 * 1024 * 1024 // 1,000 MB

// chunkSize defines the size of each chunk (in bytes) read from the input stream.
// Reading in smaller chunks helps control memory usage during large reads.
const chunkSize = 64 * 1024 // 64KB

var (
	// errReaderBytesRead indicates a failure while reading or buffering input data.
	// This error is returned when an issue occurs during the process of reading the request body.
	errReaderBytesRead = errors.New("failed to read input data")

	// errReaderLimitExceeded is returned when the read operation exceeds the maximum allowed byte limit.
	// It is triggered if the total number of bytes read surpasses the configured ReadLimit.
	errReaderLimitExceeded = errors.New("input data too large")

	// errReaderMissingBytes is returned when an attempt is made to read an empty byte slice.
	// This error occurs when the input data (Bytes) is an empty slice, indicating that there are no bytes to read.
	errReaderMissingBytes = errors.New("failed to read input bytes. Input bytes cannot be an empty slice")
)

// LimitedBytesReader is a utility for safely reading bytes with an enforced size limit.
// It is typically used to prevent reading more than a configured number of bytes
// from an incoming payload (e.g., request body).
type limitedBytesReader struct {
	// Bytes is the source data to be read.
	bytes []byte

	// ReadLimit defines the maximum number of bytes allowed to be read from Bytes.
	// If this limit is exceeded during reading, an error is returned.
	readLimit int64
}

// Read reads from the underlying byte slice up to the specified ReadLimit.
// It processes the input in 64KB chunks and returns the entire read data as a byte slice.
//
// If more than ReadLimit bytes are encountered, the function returns ErrReaderLimitExceeded.
// If the byte slice is empty, it returns ErrReaderMissingBytes.
// If an I/O or buffering error occurs during the read, it returns ErrReaderBytesRead.
func (l *limitedBytesReader) Read() ([]byte, error) {
	if len(l.bytes) == 0 {
		return nil, errReaderMissingBytes
	}

	reader := io.LimitReader(bytes.NewBuffer(l.bytes), l.readLimit+1)
	buff := bytes.NewBuffer(nil)
	bb := make([]byte, chunkSize)
	var read int64

	for {
		n, err := reader.Read(bb)
		if n > 0 {
			read += int64(n)
			if read > l.readLimit {
				return nil, errReaderLimitExceeded
			}
			_, err := buff.Write(bb[:n])
			if err != nil {
				return nil, errors.Join(err, errReaderBytesRead)
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, errors.Join(err, errReaderBytesRead)
		}
	}
	return buff.Bytes(), nil
}

// LimitOctetStreamBodyMiddleware is a Fiber middleware that limits the size of incoming
// request bodies with the Content-Type: application/octet-stream. It reads the body in chunks
// and ensures that the body does not exceed the specified size limit. The middleware responds
// with appropriate error messages if the body is empty, exceeds the size limit, or cannot be read.
//
//   - If the body exceeds the specified limit, it responds with a 400 Bad Request and an error message
//     indicating the body is too large.
//   - If the body is empty, it responds with a 400 Bad Request and an error message indicating the
//     body is missing.
//   - If an error occurs while reading the body, it responds with a 500 Internal Server Error.
//
// Usage:
//
//	app.Post("/upload", LimitOctetStreamBodyMiddleware(10*1024*1024), uploadHandler)
func LimitOctetStreamBodyMiddleware(limit int64) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !c.Is(fiber.MIMEOctetStream) {
			return c.Status(fiber.StatusBadRequest).JSON(NewUnsupportedContentTypeResponse(fiber.MIMEOctetStream))
		}

		reader := limitedBytesReader{
			bytes:     c.Body(),
			readLimit: limit,
		}

		bytes, err := reader.Read()
		switch {
		case errors.Is(err, errReaderMissingBytes):
			return c.Status(fiber.StatusBadRequest).JSON(ResponseEmptyOctetStream)

		case errors.Is(err, errReaderLimitExceeded):
			return c.Status(fiber.StatusBadRequest).JSON(NewRequestBodyTooLargeResponse(limit))

		case errors.Is(err, errReaderBytesRead):
			return c.Status(fiber.StatusInternalServerError).JSON(ResponseBodyReadFailure)

		default:
			c.Context().SetBody(bytes)
			return c.Next()
		}
	}
}

// NewRequestBodyTooLargeResponse returns a BadRequestResponse indicating that
// the request body exceeds the allowed size.
func NewRequestBodyTooLargeResponse(limit int64) openapi.BadRequestResponse {
	return openapi.BadRequestResponse{
		Message: fmt.Sprintf("The submitted octet-stream exceeds the maximum allowed size: %d.", limit),
	}
}

// NewUnsupportedContentTypeResponse returns a BadRequestResponse indicating
// that the Content-Type of the request is not supported.
func NewUnsupportedContentTypeResponse(expected string) openapi.BadRequestResponse {
	return openapi.BadRequestResponse{
		Message: fmt.Sprintf("Unsupported content type. Expected: %s.", expected),
	}
}

// ResponseBodyReadFailure is an InternalServerErrorResponse used when the server
// cannot read the request body due to an internal issue.
var ResponseBodyReadFailure = openapi.InternalServerErrorResponse{
	Message: "Failed to process request octet-stream due to a read error.",
}

// ResponseEmptyOctetStream is a BadRequestResponse used when the submitted
// request octet-stream is empty.
var ResponseEmptyOctetStream = openapi.BadRequestResponse{
	Message: "Empty request octet-stream is not allowed.",
}
