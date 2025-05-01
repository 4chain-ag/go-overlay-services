package app

import (
	"bytes"
	"errors"
	"io"
)

// ReadBodyLimit1GB defines the maximum allowed bytes readd size (in bytes).
// This limit is set to 1GB to protect against excessively large payloads.
const ReadBodyLimit1GB = 1000 * 1024 * 1024 // 1,000 MB

// chunkSize defines the size of each chunk (in bytes) read from the input stream.
// Reading in smaller chunks helps control memory usage during large reads.
const chunkSize = 64 * 1024 // 64KB

var (
	// ErrReaderBytesRead indicates a failure while reading or buffering input data.
	// This error is returned when an issue occurs during the process of reading the request body.
	ErrReaderBytesRead = errors.New("failed to read input data")

	// ErrReaderLimitExceeded is returned when the read operation exceeds the maximum allowed byte limit.
	// It is triggered if the total number of bytes read surpasses the configured ReadLimit.
	ErrReaderLimitExceeded = errors.New("input data too large")

	// ErrReaderMissingBytes is returned when an attempt is made to read an empty byte slice.
	// This error occurs when the input data (Bytes) is an empty slice, indicating that there are no bytes to read.
	ErrReaderMissingBytes = errors.New("failed to read input bytes. Input bytes cannot be an empty slice")
)

// LimitedBytesReader is a utility for safely reading bytes with an enforced size limit.
// It is typically used to prevent reading more than a configured number of bytes
// from an incoming payload (e.g., request body).
type LimitedBytesReader struct {
	// Bytes is the source data to be read.
	Bytes []byte

	// ReadLimit defines the maximum number of bytes allowed to be read from Bytes.
	// If this limit is exceeded during reading, an error is returned.
	ReadLimit int64
}

// Read reads from the underlying byte slice up to the specified ReadLimit.
// It processes the input in 64KB chunks and returns the entire read data as a byte slice.
//
// If more than ReadLimit bytes are encountered, the function returns ErrReaderLimitExceeded.
// If the byte slice is empty, it returns ErrReaderMissingBytes.
// If an I/O or buffering error occurs during the read, it returns ErrReaderBytesRead.
func (l *LimitedBytesReader) Read() ([]byte, error) {
	if len(l.Bytes) == 0 {
		return nil, ErrReaderMissingBytes
	}

	reader := io.LimitReader(bytes.NewBuffer(l.Bytes), l.ReadLimit+1)
	buff := bytes.NewBuffer(nil)
	bb := make([]byte, chunkSize)
	var read int64

	for {
		n, err := reader.Read(bb)
		if n > 0 {
			read += int64(n)
			if read > l.ReadLimit {
				return nil, ErrReaderLimitExceeded
			}
			_, err := buff.Write(bb[:n])
			if err != nil {
				return nil, errors.Join(err, ErrReaderBytesRead)
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, errors.Join(err, ErrReaderBytesRead)
		}
	}
	return buff.Bytes(), nil
}
