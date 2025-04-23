package engine

import (
	"context"
	"net/http"
	"time"

	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
)

// LookupResolver wraps the underlying lookup.LookupResolver to expose
// a simplified interface for querying and managing SLAP trackers.
type LookupResolver struct {
	resolver *lookup.LookupResolver
}

// NewLookupResolver creates and initializes a LookupResolver with a default HTTPS facilitator.
func NewLookupResolver() *LookupResolver {
	return &LookupResolver{
		resolver: &lookup.LookupResolver{
			Facilitator: &lookup.HTTPSOverlayLookupFacilitator{
				Client: http.DefaultClient,
			},
		},
	}
}

// SetSLAPTrackers configures the SLAP trackers for the resolver.
// If the given slice is empty, it leaves the resolver unchanged.
func (l *LookupResolver) SetSLAPTrackers(trackers []string) {
	if len(trackers) == 0 {
		return
	}
	l.resolver.SLAPTrackers = trackers
}

// SLAPTrackers returns the currently configured SLAP trackers.
func (l *LookupResolver) SLAPTrackers() []string {
	return l.resolver.SLAPTrackers
}

// Query performs a lookup using the configured resolver with the given question and timeout.
func (l *LookupResolver) Query(ctx context.Context, question *lookup.LookupQuestion, timeout time.Duration) (*lookup.LookupAnswer, error) {
	return l.resolver.Query(ctx, question, timeout)
}
