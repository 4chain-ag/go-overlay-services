package engine_test

import (
	"context"
	"testing"
	"time"

	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/stretchr/testify/require"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
)

func TestEngine_StartGASPSync_Success(t *testing.T) {
	t.Parallel()

	e := engine.NewEngine(engine.Engine{
		SyncConfiguration: map[string]engine.SyncConfiguration{
			"test-topic": {Type: engine.SyncConfigurationSHIP},
		},
		Advertiser: fakeAdvertiser{},
		HostingURL: "http://localhost",
	})
	e.SHIPTrackers = []string{"http://localhost"}

	e.LookupServices = map[string]engine.LookupService{
		"ls_ship": fakeLookupResolver{
			lookupAnswer: &lookup.LookupAnswer{
				Type: lookup.AnswerTypeOutputList,
				Outputs: []*lookup.OutputListItem{
					{Beef: []byte{0x00}, OutputIndex: 0},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// when
	err := e.StartGASPSync(ctx)

	// then
	require.NoError(t, err)
}


