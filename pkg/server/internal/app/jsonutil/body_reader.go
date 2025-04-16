package jsonutil

import (
	"bytes"
	"errors"
	"io"
)

// RequestBodyLimit1GB defines the maximum allowed size for request bodies (1GB).
const RequestBodyLimit1GB = 1000 * 1024 * 1024

var (
	// ErrRequestBodyRead is returned when there's an error reading the request body.
	ErrRequestBodyRead = errors.New("failed to read request body")

	// ErrRequestBodyTooLarge is returned when the request body exceeds the size limit.
	ErrRequestBodyTooLarge = errors.New("request body too large")
)

type BodyReader struct {
	Body      io.Reader
	ReadLimit int64
}

func (b *BodyReader) Read() ([]byte, error) {
	reader := io.LimitReader(b.Body, b.ReadLimit+1)
	bb := make([]byte, 64*1024)
	var buff bytes.Buffer
	var read int64

	for {
		n, err := reader.Read(bb)
		if n > 0 {
			read += int64(n)
			if read > b.ReadLimit {
				return nil, ErrRequestBodyTooLarge
			}

			if _, inner := buff.Write(bb[:n]); inner != nil {
				return nil, ErrRequestBodyRead
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, ErrRequestBodyRead
		}
	}

	return bb, nil
}
