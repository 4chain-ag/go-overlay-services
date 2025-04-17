package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/advertiser"
	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/bsv-blockchain/go-sdk/overlay/topic"
	"github.com/bsv-blockchain/go-sdk/spv"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/chaintracker"
	"golang.org/x/exp/slices"
)

var TRUE = true
var FALSE = false

type SumbitMode string

var (
	SubmitModeHistorical SumbitMode = "historical-tx"
	SubmitModeCurrent    SumbitMode = "current-tx"
)

type SyncConfigurationType int

const (
	SyncConfigurationPeers SyncConfigurationType = iota
	SyncConfigurationSHIP
	SyncConfigurationNone
)

type SyncConfiguration struct {
	Type        SyncConfigurationType
	Peers       []string
	Concurrency int
}

type OnSteakReady func(steak *overlay.Steak)
type Engine struct {
	Managers                map[string]TopicManager
	LookupServices          map[string]LookupService
	Storage                 Storage
	ChainTracker            chaintracker.ChainTracker
	HostingURL              string
	SHIPTrackers            []string
	SLAPTrackers            []string
	Broadcaster             transaction.Broadcaster
	Advertiser              advertiser.Advertiser
	SyncConfiguration       map[string]SyncConfiguration
	LogTime                 bool
	LogPrefix               string
	ErrorOnBroadcastFailure bool
	BroadcastFacilitator    topic.Facilitator
	Verbose                 bool
	PanicOnError            bool
	// Logger				  Logger //TODO: Implement Logger Interface
}

func NewEngine(cfg Engine) *Engine {
	if cfg.SyncConfiguration == nil {
		cfg.SyncConfiguration = make(map[string]SyncConfiguration)
	}
	if cfg.Managers == nil {
		cfg.Managers = make(map[string]TopicManager)
	}
	if cfg.LookupServices == nil {
		cfg.LookupServices = make(map[string]LookupService)
	}
	for name, manager := range cfg.Managers {
		config := cfg.SyncConfiguration[name]

		if name == "tm_ship" && len(cfg.SHIPTrackers) > 0 && manager != nil && config.Type == SyncConfigurationPeers {
			combined := make(map[string]struct{}, len(cfg.SHIPTrackers)+len(config.Peers))
			for _, peer := range cfg.SHIPTrackers {
				combined[peer] = struct{}{}
			}
			for _, peer := range config.Peers {
				combined[peer] = struct{}{}
			}
			config.Peers = make([]string, 0, len(combined))
			for peer := range combined {
				config.Peers = append(config.Peers, peer)
			}
			cfg.SyncConfiguration[name] = config
		} else if name == "tm_slap" && len(cfg.SLAPTrackers) > 0 && manager != nil && config.Type == SyncConfigurationPeers {
			combined := make(map[string]struct{}, len(cfg.SHIPTrackers)+len(config.Peers))
			for _, peer := range cfg.SLAPTrackers {
				combined[peer] = struct{}{}
			}
			for _, peer := range config.Peers {
				combined[peer] = struct{}{}
			}
			config.Peers = make([]string, 0, len(combined))
			for peer := range combined {
				config.Peers = append(config.Peers, peer)
			}
			cfg.SyncConfiguration[name] = config
		}
	}

	return &cfg
}

var ErrUnknownTopic = errors.New("unknown-topic")
var ErrInvalidBeef = errors.New("invalid-beef")
var ErrInvalidTransaction = errors.New("invalid-transaction")
var ErrMissingInput = errors.New("missing-input")
var ErrInputSpent = errors.New("input-spent")

func (e *Engine) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode SumbitMode, onSteakReady OnSteakReady) (overlay.Steak, error) {
	start := time.Now()
	for _, topic := range taggedBEEF.Topics {
		if _, ok := e.Managers[topic]; !ok {
			return nil, ErrUnknownTopic
		}
	}

	var tx *transaction.Transaction
	beef, tx, txid, err := transaction.ParseBeef(taggedBEEF.Beef)
	if err != nil {
		if e.PanicOnError {
			log.Panicln(err)
		}
		return nil, err
	} else if tx == nil {
		if e.PanicOnError {
			log.Panicln(ErrInvalidBeef)
		}
		return nil, ErrInvalidBeef
	}
	if valid, err := spv.Verify(tx, e.ChainTracker, nil); err != nil {
		if e.PanicOnError {
			log.Panicln(err)
		}
		return nil, err
	} else if !valid {
		if e.PanicOnError {
			log.Panicln(ErrInvalidTransaction)
		}
		return nil, ErrInvalidTransaction
	}
	if e.Verbose {
		fmt.Println("Validated in", time.Since(start))
		start = time.Now()
	}
	steak := make(overlay.Steak, len(taggedBEEF.Topics))
	topicInputs := make(map[string]map[uint32]*Output, len(tx.Inputs))
	inpoints := make([]*overlay.Outpoint, 0, len(tx.Inputs))
	ancillaryBeefs := make(map[string][]byte, len(taggedBEEF.Topics))
	for _, input := range tx.Inputs {
		inpoints = append(inpoints, &overlay.Outpoint{
			Txid:        *input.SourceTXID,
			OutputIndex: input.SourceTxOutIndex,
		})
	}
	dupeTopics := make(map[string]struct{}, len(taggedBEEF.Topics))
	for _, topic := range taggedBEEF.Topics {
		if exists, err := e.Storage.DoesAppliedTransactionExist(ctx, &overlay.AppliedTransaction{
			Txid:  txid,
			Topic: topic,
		}); err != nil {
			if e.PanicOnError {
				log.Panicln(err)
			}
			return nil, err
		} else if exists {
			steak[topic] = &overlay.AdmittanceInstructions{}
			dupeTopics[topic] = struct{}{}
			continue
		} else {
			topicInputs[topic] = make(map[uint32]*Output, len(tx.Inputs))
			previousCoins := make(map[uint32][]byte, len(tx.Inputs))
			for vin, outpoint := range inpoints {
				if output, err := e.Storage.FindOutput(ctx, outpoint, &topic, nil, true); err != nil {
					return nil, err
				} else if output != nil {
					previousCoins[uint32(vin)] = output.Beef
					topicInputs[topic][uint32(vin)] = output
				}
			}

			if admit, err := e.Managers[topic].IdentifyAdmissableOutputs(ctx, taggedBEEF.Beef, previousCoins); err != nil {
				if e.PanicOnError {
					log.Panicln(err)
				}
				return nil, err
			} else {
				if e.Verbose {
					fmt.Println("Identified in", time.Since(start))
					start = time.Now()
				}
				if len(admit.AncillaryTxids) > 0 {
					ancillaryBeef := transaction.Beef{
						Version:      transaction.BEEF_V2,
						Transactions: make(map[string]*transaction.BeefTx, len(admit.AncillaryTxids)),
					}
					for _, txid := range admit.AncillaryTxids {
						if tx := beef.FindTransaction(txid.String()); tx == nil {
							return nil, errors.New("missing dependency transaction")
						} else if beefBytes, err := tx.BEEF(); err != nil {
							return nil, err
						} else if err := ancillaryBeef.MergeBeefBytes(beefBytes); err != nil {
							return nil, err
						}
					}
					if beefBytes, err := ancillaryBeef.Bytes(); err != nil {
						return nil, err
					} else {
						ancillaryBeefs[topic] = beefBytes
					}
				}
				steak[topic] = &admit
			}
		}
	}

	for _, topic := range taggedBEEF.Topics {
		if _, ok := dupeTopics[topic]; ok {
			continue
		}
		for _, outpoint := range inpoints {
			if err := e.Storage.MarkUTXOAsSpent(ctx, outpoint, topic); err != nil {
				if e.PanicOnError {
					log.Panicln(err)
				}
				return nil, err
			} else {
				for _, l := range e.LookupServices {
					for _, inpoint := range inpoints {
						if err := l.OutputSpent(ctx, inpoint, topic, taggedBEEF.Beef); err != nil {
							if e.PanicOnError {
								log.Panicln(err)
							}
							return nil, err
						}
					}
				}
			}
		}
	}
	if e.Verbose {
		fmt.Println("Marked spent in", time.Since(start))
		start = time.Now()
	}
	if mode != SubmitModeHistorical && e.Broadcaster != nil {
		if _, failure := e.Broadcaster.Broadcast(tx); failure != nil {
			if e.PanicOnError {
				log.Panicln(err)
			}
			return nil, failure
		}
	}

	if onSteakReady != nil {
		onSteakReady(&steak)
	}

	for _, topic := range taggedBEEF.Topics {
		if _, ok := dupeTopics[topic]; ok {
			continue
		}
		admit := steak[topic]
		outputsConsumed := make([]*Output, 0, len(admit.CoinsToRetain))
		outpointsConsumed := make([]*overlay.Outpoint, 0, len(admit.CoinsToRetain))
		for vin, output := range topicInputs[topic] {
			for _, coin := range admit.CoinsToRetain {
				if vin == coin {
					outputsConsumed = append(outputsConsumed, output)
					outpointsConsumed = append(outpointsConsumed, &output.Outpoint)
					delete(topicInputs[topic], vin)
					break
				}
			}
		}

		for vin, output := range topicInputs[topic] {
			if err := e.deleteUTXODeep(ctx, output); err != nil {
				if e.PanicOnError {
					log.Panicln(err)
				}
				return nil, err
			}
			admit.CoinsRemoved = append(admit.CoinsRemoved, uint32(vin))
		}

		newOutpoints := make([]*overlay.Outpoint, 0, len(admit.OutputsToAdmit))
		for _, vout := range admit.OutputsToAdmit {
			out := tx.Outputs[vout]
			output := &Output{
				Outpoint: overlay.Outpoint{
					Txid:        *txid,
					OutputIndex: uint32(vout),
				},
				Script:          out.LockingScript,
				Satoshis:        out.Satoshis,
				Topic:           topic,
				OutputsConsumed: outpointsConsumed,
				Beef:            taggedBEEF.Beef,
				AncillaryTxids:  admit.AncillaryTxids,
				AncillaryBeef:   ancillaryBeefs[topic],
			}
			if tx.MerklePath != nil {
				output.BlockHeight = tx.MerklePath.BlockHeight
				for _, leaf := range tx.MerklePath.Path[0] {
					if leaf.Hash != nil && leaf.Hash.Equal(output.Outpoint.Txid) {
						output.BlockIdx = leaf.Offset
						break
					}
				}
			}
			if err := e.Storage.InsertOutput(ctx, output); err != nil {
				if e.PanicOnError {
					log.Panicln(err)
				}
				return nil, err
			}
			newOutpoints = append(newOutpoints, &output.Outpoint)
			for _, l := range e.LookupServices {
				if err := l.OutputAdded(ctx, &output.Outpoint, topic, output.Beef); err != nil {
					if e.PanicOnError {
						log.Panicln(err)
					}
					return nil, err
				}
			}
		}
		if e.Verbose {
			fmt.Println("Outputs added in", time.Since(start))
			start = time.Now()
		}
		for _, output := range outputsConsumed {
			output.ConsumedBy = append(output.ConsumedBy, newOutpoints...)

			if err := e.Storage.UpdateConsumedBy(ctx, &output.Outpoint, output.Topic, output.ConsumedBy); err != nil {
				if e.PanicOnError {
					log.Panicln(err)
				}
				return nil, err
			}
		}
		if e.Verbose {
			fmt.Println("Consumes updated in ", time.Since(start))
			start = time.Now()
		}
		if err := e.Storage.InsertAppliedTransaction(ctx, &overlay.AppliedTransaction{
			Txid:  txid,
			Topic: topic,
		}); err != nil {
			if e.PanicOnError {
				log.Panicln(err)
			}
			return nil, err
		}
		if e.Verbose {
			fmt.Println("Applied in", time.Since(start))
		}
	}
	if e.Advertiser == nil || mode == SubmitModeHistorical {
		return steak, nil
	}

	releventTopics := make([]string, 0, len(taggedBEEF.Topics))
	for topic, steak := range steak {
		if steak.OutputsToAdmit == nil && steak.CoinsToRetain == nil {
			continue
		}
		if _, ok := dupeTopics[topic]; !ok {
			releventTopics = append(releventTopics, topic)
		}
	}
	if len(releventTopics) == 0 {
		return steak, nil
	}

	broadcasterCfg := &topic.BroadcasterConfig{}
	if len(e.SLAPTrackers) > 0 {
		broadcasterCfg.Resolver = lookup.NewLookupResolver(&lookup.LookupResolver{
			SLAPTrackers: e.SLAPTrackers,
		})
	}

	if broadcaster, err := topic.NewBroadcaster(releventTopics, broadcasterCfg); err != nil {
		log.Println("Error during propagation to other nodes:", err)
	} else if _, failure := broadcaster.BroadcastCtx(ctx, tx); failure != nil {
		log.Println("Error during propagation to other nodes:", failure)
	}
	return steak, nil
}

func (e *Engine) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	if l, ok := e.LookupServices[question.Service]; !ok {
		return nil, ErrUnknownTopic
	} else if result, err := l.Lookup(ctx, question); err != nil {
		return nil, err
	} else if result.Type == lookup.AnswerTypeFreeform || result.Type == lookup.AnswerTypeOutputList {
		return result, nil
	} else {
		hydratedOutputs := make([]*lookup.OutputListItem, 0, len(result.Outputs))
		for _, formula := range result.Formulas {
			if output, err := e.Storage.FindOutput(ctx, formula.Outpoint, nil, nil, true); err != nil {
				return nil, err
			} else if output != nil && output.Beef != nil {
				if output, err := e.GetUTXOHistory(ctx, output, formula.Histoy, 0); err != nil {
					return nil, err
				} else if output != nil {
					hydratedOutputs = append(hydratedOutputs, &lookup.OutputListItem{
						Beef:        output.Beef,
						OutputIndex: output.Outpoint.OutputIndex,
					})
				}
			}
		}
		return &lookup.LookupAnswer{
			Type:    lookup.AnswerTypeOutputList,
			Outputs: hydratedOutputs,
		}, nil
	}
}

func (e *Engine) GetUTXOHistory(ctx context.Context, output *Output, historySelector func(beef []byte, outputIndex uint32, currentDepth uint32) bool, currentDepth uint32) (*Output, error) {
	if historySelector == nil {
		return output, nil
	}
	shouldTravelHistory := historySelector(output.Beef, output.Outpoint.OutputIndex, currentDepth)
	if !shouldTravelHistory {
		return nil, nil
	}
	if output != nil && len(output.OutputsConsumed) == 0 {
		return output, nil
	}
	outputsConsumed := output.OutputsConsumed[:]
	childHistories := make(map[string]*Output, len(outputsConsumed))
	for _, outpoint := range outputsConsumed {
		if output, err := e.Storage.FindOutput(ctx, outpoint, nil, nil, true); err != nil {
			return nil, err
		} else if output != nil {
			if child, err := e.GetUTXOHistory(ctx, output, historySelector, currentDepth+1); err != nil {
				return nil, err
			} else if child != nil {
				childHistories[child.Outpoint.String()] = child
			}
		}
	}

	if tx, err := transaction.NewTransactionFromBEEF(output.Beef); err != nil {
		return nil, err
	} else {
		for _, txin := range tx.Inputs {
			outpoint := &overlay.Outpoint{
				Txid:        *txin.SourceTXID,
				OutputIndex: txin.SourceTxOutIndex,
			}
			if input := childHistories[outpoint.String()]; input != nil {
				if input.Beef == nil {
					return nil, errors.New("missing beef")
				} else if txin.SourceTransaction, err = transaction.NewTransactionFromBEEF(input.Beef); err != nil {
					return nil, err
				}
			}
		}
		if beef, err := tx.BEEF(); err != nil {
			return nil, err
		} else {
			output.Beef = beef
			return output, nil
		}
	}
}

func (e *Engine) SyncAdvertisements(ctx context.Context) error {
	if e.Advertiser == nil {
		return nil
	}
	configuredTopics := make([]string, 0, len(e.Managers))
	requiredSHIPAdvertisements := make(map[string]struct{}, len(configuredTopics))
	for name := range e.Managers {
		configuredTopics = append(configuredTopics, name)
		requiredSHIPAdvertisements[name] = struct{}{}
	}
	configuredServices := make([]string, 0, len(e.LookupServices))
	requiredSLAPAdvertisements := make(map[string]struct{}, len(configuredServices))
	for name := range e.LookupServices {
		configuredServices = append(configuredServices, name)
		requiredSLAPAdvertisements[name] = struct{}{}
	}
	currentSHIPAdvertisements, err := e.Advertiser.FindAllAdvertisements("SHIP")
	if err != nil {
		return err
	}
	shipsToCreate := make([]string, 0, len(requiredSHIPAdvertisements))
	for topic := range requiredSHIPAdvertisements {
		if slices.IndexFunc(currentSHIPAdvertisements, func(ad *advertiser.Advertisement) bool {
			return ad.TopicOrService == topic && ad.Domain == e.HostingURL
		}) == -1 {
			shipsToCreate = append(shipsToCreate, topic)
		}
	}
	shipsToRevoke := make([]*advertiser.Advertisement, 0, len(currentSHIPAdvertisements))
	for _, ad := range currentSHIPAdvertisements {
		if _, ok := requiredSHIPAdvertisements[ad.TopicOrService]; !ok {
			shipsToRevoke = append(shipsToRevoke, ad)
		}
	}

	currentSLAPAdvertisements, err := e.Advertiser.FindAllAdvertisements("SLAP")
	if err != nil {
		return err
	}
	slapsToCreate := make([]string, 0, len(requiredSLAPAdvertisements))
	for service := range requiredSLAPAdvertisements {
		if slices.IndexFunc(currentSLAPAdvertisements, func(ad *advertiser.Advertisement) bool {
			return ad.TopicOrService == service && ad.Domain == e.HostingURL
		}) == -1 {
			slapsToCreate = append(slapsToCreate, service)
		}
	}
	slapsToRevoke := make([]*advertiser.Advertisement, 0, len(currentSLAPAdvertisements))
	for _, ad := range currentSLAPAdvertisements {
		if _, ok := requiredSLAPAdvertisements[ad.TopicOrService]; !ok {
			slapsToRevoke = append(slapsToRevoke, ad)
		}
	}
	advertisementData := make([]*advertiser.AdvertisementData, 0, len(shipsToCreate)+len(slapsToCreate))
	for _, topic := range shipsToCreate {
		advertisementData = append(advertisementData, &advertiser.AdvertisementData{
			Protocol:           "SHIP",
			TopicOrServiceName: topic,
		})
	}
	for _, service := range slapsToCreate {
		advertisementData = append(advertisementData, &advertiser.AdvertisementData{
			Protocol:           "SLAP",
			TopicOrServiceName: service,
		})
	}
	if len(advertisementData) > 0 {
		if taggedBEEF, err := e.Advertiser.CreateAdvertisements(advertisementData); err != nil {
			log.Println("Failed to create SHIP advertisement:", err)
		} else if _, err := e.Submit(ctx, taggedBEEF, SubmitModeCurrent, nil); err != nil {
			log.Println("Failed to create SHIP advertisement:", err)
		}
	}
	revokeData := make([]*advertiser.Advertisement, 0, len(shipsToRevoke)+len(slapsToRevoke))
	revokeData = append(revokeData, shipsToRevoke...)
	revokeData = append(revokeData, slapsToRevoke...)
	if len(revokeData) > 0 {
		if taggedBEEF, err := e.Advertiser.RevokeAdvertisements(revokeData); err != nil {
			log.Println("Failed to revoke SHIP/SLAP advertisements:", err)
		} else if _, err := e.Submit(ctx, taggedBEEF, SubmitModeCurrent, nil); err != nil {
			log.Println("Failed to revoke SHIP/SLAP advertisements:", err)
		}
	}
	return nil
}

func (e *Engine) StartGASPSync(ctx context.Context) error {
	if e.SyncConfiguration == nil {
		return errors.New("not configured for topical synchronization")
	}

	for topic := range e.SyncConfiguration {
		syncEndpoints, ok := e.SyncConfiguration[topic]
		if !ok {
			continue
		}
		if syncEndpoints.Type == SyncConfigurationSHIP {
			resolver := lookup.LookupResolver{
				Facilitator: &lookup.HTTPSOverlayLookupFacilitator{
					Client: http.DefaultClient,
				},
			}
			if e.SLAPTrackers != nil {
				resolver.SLAPTrackers = e.SLAPTrackers
			}

			if query, err := json.Marshal(map[string]any{
				"topics": []string{topic},
			}); err != nil {
				return err
			} else if lookupAnswer, err := resolver.Query(ctx, &lookup.LookupQuestion{
				Service: "ls_ship",
				Query:   query,
			}, 10*time.Second); err != nil {
				return err
			} else if lookupAnswer.Type == lookup.AnswerTypeOutputList {
				endpointSet := make(map[string]struct{}, len(lookupAnswer.Outputs))
				for _, output := range lookupAnswer.Outputs {
					if tx, err := transaction.NewTransactionFromBEEF(output.Beef); err != nil {
						log.Println("Failed to parse advertisement output:", err)
					} else if advertisement, err := e.Advertiser.ParseAdvertisement(tx.Outputs[output.OutputIndex].LockingScript); err != nil {
						log.Println("Failed to parse advertisement output:", err)
					} else if advertisement != nil && advertisement.Protocol == "SHIP" {
						endpointSet[advertisement.Domain] = struct{}{}
					}
				}
				syncEndpoints.Peers = make([]string, 0, len(endpointSet))
				for endpoint := range endpointSet {
					if endpoint != e.HostingURL {
						syncEndpoints.Peers = append(syncEndpoints.Peers, endpoint)
					}
				}
			}
		}

		if len(syncEndpoints.Peers) > 0 {
			peers := make([]string, 0, len(syncEndpoints.Peers))
			for _, peer := range syncEndpoints.Peers {
				if peer != e.HostingURL {
					peers = append(peers, peer)
				}
			}
			for _, peer := range peers {
				logPrefix := "[GASP Sync of " + topic + " with " + peer + "]"
				gasp := core.NewGASP(core.GASPParams{
					Storage: NewOverlayGASPStorage(topic, e, nil),
					Remote: &OverlayGASPRemote{
						EndpointUrl: peer,
						Topic:       topic,
						HttpClient:  http.DefaultClient,
					},
					LogPrefix:      &logPrefix,
					Unidirectional: true,
					Concurrency:    syncEndpoints.Concurrency,
				})
				if err := gasp.Sync(ctx); err != nil {
					log.Println("Failed to sync with peer", peer, ":", err)
				}
			}
		}
	}
	return nil
}

func (e *Engine) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	if utxos, err := e.Storage.FindUTXOsForTopic(ctx, topic, initialRequest.Since, false); err != nil {
		return nil, err
	} else {
		utxoList := make([]*overlay.Outpoint, 0, len(utxos))
		for _, utxo := range utxos {
			utxoList = append(utxoList, &utxo.Outpoint)
		}
		return &core.GASPInitialResponse{
			UTXOList: utxoList,
		}, nil
	}
}

func (e *Engine) ProvideForeignGASPNode(ctx context.Context, graphId *overlay.Outpoint, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	var hydrator func(ctx context.Context, output *Output) (*core.GASPNode, error)
	hydrator = func(ctx context.Context, output *Output) (*core.GASPNode, error) {
		if output.Beef == nil {
			return nil, ErrMissingInput
		} else if _, tx, _, err := transaction.ParseBeef(output.Beef); err != nil {
			return nil, err
		} else if tx == nil {
			for _, outpoint := range output.OutputsConsumed {
				if output, err := e.Storage.FindOutput(ctx, outpoint, &topic, nil, false); err == nil {
					return hydrator(ctx, output)
				}
			}
			return nil, errors.New("unable to find output")
		} else {
			node := &core.GASPNode{
				GraphID:       graphId,
				RawTx:         tx.Hex(),
				OutputIndex:   outpoint.OutputIndex,
				AncillaryBeef: output.AncillaryBeef,
			}
			if tx.MerklePath != nil {
				proof := tx.MerklePath.Hex()
				node.Proof = &proof
			}
			return node, nil

		}

	}
	if output, err := e.Storage.FindOutput(ctx, graphId, &topic, nil, true); err != nil {
		return nil, err
	} else {
		return hydrator(ctx, output)
	}
}

func (e *Engine) deleteUTXODeep(ctx context.Context, output *Output) error {
	if len(output.ConsumedBy) == 0 {
		if err := e.Storage.DeleteOutput(ctx, &output.Outpoint, output.Topic); err != nil {
			return err
		}
		for _, l := range e.LookupServices {
			if err := l.OutputDeleted(ctx, &output.Outpoint, output.Topic); err != nil {
				return err
			}
		}
	}
	if len(output.OutputsConsumed) == 0 {
		return nil
	}

	for _, outpoint := range output.OutputsConsumed {
		staleOutput, err := e.Storage.FindOutput(ctx, outpoint, &output.Topic, nil, false)
		if err != nil {
			return err
		} else if staleOutput == nil {
			continue
		}
		if len(staleOutput.ConsumedBy) > 0 {
			consumedBy := staleOutput.ConsumedBy
			staleOutput.ConsumedBy = make([]*overlay.Outpoint, 0, len(consumedBy))
			for _, outpoint := range consumedBy {
				if !bytes.Equal(outpoint.TxBytes(), output.Outpoint.TxBytes()) {
					staleOutput.ConsumedBy = append(staleOutput.ConsumedBy, outpoint)
				}
			}
			if err := e.Storage.UpdateConsumedBy(ctx, &staleOutput.Outpoint, staleOutput.Topic, staleOutput.ConsumedBy); err != nil {
				return err
			}
		}

		if err := e.deleteUTXODeep(ctx, staleOutput); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) updateInputProofs(ctx context.Context, tx *transaction.Transaction, txid chainhash.Hash, proof *transaction.MerklePath) (err error) {
	if tx.MerklePath != nil {
		tx.MerklePath = proof
		return
	}

	if tx.TxID().Equal(txid) {
		tx.MerklePath = proof
	} else {
		for _, input := range tx.Inputs {
			if input.SourceTransaction == nil {
				return errors.New("missing source transaction")
			} else if err = e.updateInputProofs(ctx, input.SourceTransaction, txid, proof); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Engine) updateMerkleProof(ctx context.Context, output *Output, txid chainhash.Hash, proof *transaction.MerklePath) error {
	if len(output.Beef) == 0 {
		return errors.New("missing beef")
	}
	beef, tx, _, err := transaction.ParseBeef(output.Beef)
	if err != nil {
		return err
	} else if tx == nil {
		return errors.New("missing transaction")
	}
	if tx.MerklePath != nil {
		if oldRoot, err := tx.MerklePath.ComputeRoot(&txid); err != nil {
			return err
		} else if newRoot, err := proof.ComputeRoot(&txid); err != nil {
			return err
		} else if oldRoot.Equal(*newRoot) {
			return nil
		}
	}
	if err = e.updateInputProofs(ctx, tx, txid, proof); err != nil {
		return err
	} else if atomicBytes, err := tx.AtomicBEEF(false); err != nil {
		return err
	} else {
		if len(output.AncillaryTxids) > 0 {
			ancillaryBeef := transaction.Beef{
				Version:      transaction.BEEF_V2,
				Transactions: make(map[string]*transaction.BeefTx, len(output.AncillaryTxids)),
			}
			for _, dep := range output.AncillaryTxids {
				if depTx := beef.FindTransaction(dep.String()); depTx == nil {
					return errors.New("missing dependency transaction")
				} else if depBeefBytes, err := depTx.BEEF(); err != nil {
					return err
				} else if err := ancillaryBeef.MergeBeefBytes(depBeefBytes); err != nil {
					return err
				}
			}
			if output.AncillaryBeef, err = ancillaryBeef.Bytes(); err != nil {
				return err
			}
		} else {
			output.AncillaryBeef = nil
		}

		output.BlockHeight = proof.BlockHeight
		for _, leaf := range proof.Path[0] {
			if leaf.Hash != nil && leaf.Hash.Equal(output.Outpoint.Txid) {
				output.BlockIdx = leaf.Offset
				break
			}
		}
		if err = e.Storage.UpdateTransactionBEEF(ctx, &output.Outpoint.Txid, atomicBytes); err != nil {
			return err
		}
		for _, outpoint := range output.ConsumedBy {
			if consumingOutputs, err := e.Storage.FindOutputsForTransaction(ctx, &outpoint.Txid, true); err != nil {
				return err
			} else {
				for _, consuming := range consumingOutputs {
					if err := e.updateMerkleProof(ctx, consuming, txid, proof); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil

}

func (e *Engine) HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error {
	if outputs, err := e.Storage.FindOutputsForTransaction(ctx, txid, true); err != nil {
		return err
	} else {
		for _, output := range outputs {
			if err := e.updateMerkleProof(ctx, output, *txid, proof); err != nil {
				return err
			} else if err := e.Storage.UpdateOutputBlockHeight(ctx, &output.Outpoint, output.Topic, output.BlockHeight, output.BlockIdx, output.AncillaryBeef); err != nil {
				return err
			}
			for _, l := range e.LookupServices {
				if err := l.OutputBlockHeightUpdated(ctx, &output.Outpoint, output.BlockHeight, output.BlockIdx); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (e *Engine) ListTopicManagers() map[string]*overlay.MetaData {
	result := make(map[string]*overlay.MetaData, len(e.Managers))
	for name, manager := range e.Managers {
		result[name] = manager.GetMetaData()
	}
	return result
}

func (e *Engine) ListLookupServiceProviders() map[string]*overlay.MetaData {
	result := make(map[string]*overlay.MetaData, len(e.LookupServices))
	for name, provider := range e.LookupServices {
		result[name] = provider.GetMetaData()
	}
	return result
}

func (e *Engine) GetDocumentationForTopicManager(manager string) (string, error) {
	if tm, ok := e.Managers[manager]; !ok {
		return "", errors.New("no documentation found")
	} else {
		return tm.GetDocumentation(), nil
	}
}

func (e *Engine) GetDocumentationForLookupServiceProvider(provider string) (string, error) {
	if l, ok := e.LookupServices[provider]; !ok {
		return "", errors.New("no documentation found")
	} else {
		return l.GetDocumentation(), nil
	}
}
