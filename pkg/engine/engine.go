package engine

import "context"

// NoopEngineProvider is a custom test overlay engine implementation. This is only a temporary solution and will be removed
// after migrating the engine code. Currently, it functions as mock for the overlay HTTP server.
type NoopEngineProvider struct{}

// SubmitTransaction is a no-op call that always returns a nil error.
func (*NoopEngineProvider) SubmitTransaction(ctx context.Context) error { return nil }

// SyncAdvertisments is a no-op call that always returns a nil error.
func (*NoopEngineProvider) SyncAdvertisments(ctx context.Context) error { return nil }

// GetTopicManagerDocumentation is a no-op call that always returns a nil error.
func (*NoopEngineProvider) GetTopicManagerDocumentation(ctx context.Context) error { return nil }
