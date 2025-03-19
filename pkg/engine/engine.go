package engine

import "context"

type NoopEngineProvider struct{}

func (*NoopEngineProvider) SubmitTransaction(ctx context.Context) error { return nil }

func (*NoopEngineProvider) SyncAdvertisments(ctx context.Context) error { return nil }

func (*NoopEngineProvider) GetTopicManagerDocumentation(ctx context.Context) error { return nil }
