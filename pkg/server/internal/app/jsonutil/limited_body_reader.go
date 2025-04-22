package jsonutil

import (
	"bytes"
	"errors"
	"io"
)

// RequestBodyLimit1GB defines the maximum allowed size for request bodies (1GB).
const RequestBodyLimit1GB = 1000 * 1024 * 1024

// chunkSize defines the size of chunks to use when reading request bodies (64KB).
const chunkSize = 64 * 1024

var (
	// ErrRequestBodyRead is returned when there's an error reading the request body.
	ErrRequestBodyRead = errors.New("failed to read request body")
	// ErrRequestBodyTooLarge is returned when the request body exceeds the size limit.
	ErrRequestBodyTooLarge = errors.New("request body too large")
	// ErrBodyReaderFailure is returned when the body reader fails unexpectedly,
	// without returning ErrRequestBodyRead or ErrRequestBodyTooLarge.
	ErrBodyReaderFailure = errors.New("unexpected failure while reading request body")
)

// LimitedBodyReader is a struct that reads and processes data from an io.Reader
// with a specified limit on the maximum number of bytes that can be read.
// It helps manage reading large request bodies by ensuring the data doesn't
// exceed a predefined size limit.
type LimitedBodyReader struct {
	Body      io.Reader
	ReadLimit int64
}

// Read reads data from the LimitedBodyReader's io.Reader and returns the data as a byte slice.
// It will read chunks of data from the body and accumulate them until the read limit is reached
// or until the end of the body is encountered. If the read limit is exceeded, it returns an error.
// If the reading or writing to the buffer fails, an error is returned.
//
// The function reads the data in chunks of 64KB (defined by the buffer size),
// ensuring that no more than the allowed `ReadLimit` bytes are read.
func (l *LimitedBodyReader) Read() ([]byte, error) {
	reader := io.LimitReader(l.Body, l.ReadLimit+1)
	buff := bytes.NewBuffer(nil)
	bb := make([]byte, chunkSize)
	var read int64

	for {
		n, err := reader.Read(bb)
		if n > 0 {
			read += int64(n)
			if read > l.ReadLimit {
				return nil, ErrRequestBodyTooLarge
			}
			_, err := buff.Write(bb[:n])
			if err != nil {
				return nil, ErrRequestBodyRead
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, errors.Join(err, ErrBodyReaderFailure)
		}
	}
	return buff.Bytes(), nil
}
