package app_test

import (
	"strings"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/stretchr/testify/require"
)

func TestLimitedBytesReader_NegativeScenarios(t *testing.T) {
	tests := map[string]struct {
		name        string
		bytes       []byte
		readLimit   int64
		expectError error
	}{
		"body exceeds limit": {
			bytes:       []byte(strings.Repeat("A", 1025)),
			readLimit:   1024,
			expectError: app.ErrReaderLimitExceeded,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			reader := &app.LimitedBytesReader{
				Bytes:     tc.bytes,
				ReadLimit: tc.readLimit,
			}

			// when:
			data, err := reader.Read()

			// then:
			require.ErrorIs(t, err, tc.expectError)
			require.Nil(t, data)
		})
	}
}
