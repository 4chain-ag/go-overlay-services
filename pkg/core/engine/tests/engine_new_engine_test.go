package engine_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/stretchr/testify/require"
)

func TestEngine_NewEngine_ShouldInitializeFields_WhenNilProvided(t *testing.T) {
	t.Parallel()

	// given:
	input := engine.Engine{}

	// when:
	result := engine.NewEngine(input)

	// then:
	require.NotNil(t, result)
	require.NotNil(t, result.Managers, "Managers should be initialized")
	require.NotNil(t, result.LookupServices, "LookupServices should be initialized")
	require.NotNil(t, result.SyncConfiguration, "SyncConfiguration should be initialized")
}

func TestEngine_NewEngine_ShouldMergeTrackers_WhenManagerIsShipType(t *testing.T) {
	t.Parallel()

	// given:
	input := engine.Engine{
		SHIPTrackers: []string{"http://tracker1.com"},
		Managers: map[string]engine.TopicManager{
			"tm_ship": fakeTopicManager{},
		},
		SyncConfiguration: map[string]engine.SyncConfiguration{
			"tm_ship": {Type: engine.SyncConfigurationPeers, Peers: []string{"http://peer1.com"}},
		},
	}

	// when:
	result := engine.NewEngine(input)

	// then:
	require.NotNil(t, result)
	peers := result.SyncConfiguration["tm_ship"].Peers
	require.ElementsMatch(t, peers, []string{"http://tracker1.com", "http://peer1.com"})
}

func TestEngine_NewEngine_ShouldMergeTrackers_WhenManagerIsSlapType(t *testing.T) {
	t.Parallel()

	// given:
	input := engine.Engine{
		SLAPTrackers: []string{"http://slaptracker.com"},
		Managers: map[string]engine.TopicManager{
			"tm_slap": fakeTopicManager{},
		},
		SyncConfiguration: map[string]engine.SyncConfiguration{
			"tm_slap": {Type: engine.SyncConfigurationPeers, Peers: []string{"http://peer2.com"}},
		},
	}

	// when:
	result := engine.NewEngine(input)

	// then:
	require.NotNil(t, result)
	peers := result.SyncConfiguration["tm_slap"].Peers
	require.ElementsMatch(t, peers, []string{"http://slaptracker.com", "http://peer2.com"})
}

func TestEngine_NewEngine_ShouldNotMergeTrackers_WhenTypeIsNotPeers(t *testing.T) {
	t.Parallel()

	// given:
	input := engine.Engine{
		SHIPTrackers: []string{"http://tracker-should-not-merge.com"},
		Managers: map[string]engine.TopicManager{
			"tm_ship": fakeTopicManager{},
		},
		SyncConfiguration: map[string]engine.SyncConfiguration{
			"tm_ship": {Type: engine.SyncConfigurationSHIP, Peers: []string{"http://peer1.com"}},
		},
	}

	// when:
	result := engine.NewEngine(input)

	// then:
	require.NotNil(t, result)
	peers := result.SyncConfiguration["tm_ship"].Peers
	require.ElementsMatch(t, peers, []string{"http://peer1.com"}, "Trackers should not be merged if type != Peers")
}
