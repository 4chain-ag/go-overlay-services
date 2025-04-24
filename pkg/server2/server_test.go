package server2_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/stretchr/testify/require"
)

func Test_HTTPServer_ShouldShutdownAfterContextTimeout(t *testing.T) {
	// given:
	srv := server2.New()

	// when:
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	done := srv.ListenAndServe(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// then:
		_, ok := <-done
		require.False(t, ok, "Server did not shut down after context timeout")
	}()

	wg.Wait()
}
