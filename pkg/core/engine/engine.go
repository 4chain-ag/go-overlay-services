package engine

import (
	"bytes"
	"context"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/core/advertiser"
	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/bsv-blockchain/go-sdk/overlay/topic"
	"github.com/bsv-blockchain/go-sdk/spv"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/chaintracker"
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
	Type  SyncConfigurationType
	Peers []string
}

type OnSteakReady func(steak overlay.Steak)
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

func (e *Engine) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode SumbitMode, onSteakReady OnSteakReady) (overlay.Steak, error) {
	for _, topic := range taggedBEEF.Topics {
		if _, ok := e.Managers[topic]; !ok {
			return nil, ErrUnknownTopic
		}
	}

	if beef, txid, err := transaction.NewBeefFromAtomicBytes(taggedBEEF.Beef); err != nil {
		return nil, ErrInvalidBeef
	} else if tx := beef.FindTransaction(txid.String()); tx == nil {
		return nil, ErrInvalidBeef
	} else if valid, err := spv.Verify(tx, e.ChainTracker, nil); err != nil {
		return nil, err
	} else if !valid {
		return nil, ErrInvalidTransaction
	} else {
		steak := make(overlay.Steak, len(taggedBEEF.Topics))
		tx := beef.FindTransaction(txid.String())
		topicInputs := make(map[string]map[uint32]*Output, len(tx.Inputs))
		inpoints := make([]*overlay.Outpoint, 0, len(tx.Inputs))
		for _, input := range tx.Inputs {
			inpoints = append(inpoints, &overlay.Outpoint{
				Txid:        *input.SourceTXID,
				OutputIndex: input.SourceTxOutIndex,
			})
		}
		for _, topic := range taggedBEEF.Topics {
			if exists, err := e.Storage.DoesAppliedTransactionExist(ctx, &overlay.AppliedTransaction{
				Txid:  txid,
				Topic: topic,
			}); err != nil {
				return nil, err
			} else if exists {
				steak[topic] = &overlay.AdmittanceInstructions{}
				continue
			} else {
				previousCoins := make([]uint32, 0, len(tx.Inputs))
				if inputs, err := e.Storage.FindOutputs(ctx, inpoints, &topic, &FALSE, false); err != nil {
					return nil, err
				} else {
					topicInputs[topic] = make(map[uint32]*Output, len(inputs))
					for vin, input := range inputs {
						if input != nil {
							previousCoins = append(previousCoins, uint32(vin))
							topicInputs[topic][input.Outpoint.OutputIndex] = input
						}
					}
				}
				if admit, err := e.Managers[topic].IdentifyAdmissableOutputs(ctx, beef, txid, previousCoins); err != nil {
					return nil, err
				} else {
					steak[topic] = &admit
				}
			}
		}
		for _, topic := range taggedBEEF.Topics {
			if err := e.Storage.MarkUTXOsAsSpent(ctx, inpoints, topic); err != nil {
				return nil, err
			}
		}
		if mode != SubmitModeHistorical && e.Broadcaster != nil {
			if _, failure := e.Broadcaster.Broadcast(tx); failure != nil {
				return nil, failure
			}
		}

		if onSteakReady != nil {
			onSteakReady(steak)
		}

		for _, topic := range taggedBEEF.Topics {
			admit := steak[topic]
			outputsConsumed := make([]*Output, 0, len(admit.CoinsToRetain))
			outpointsConsumed := make([]*overlay.Outpoint, 0, len(admit.CoinsToRetain))
			for _, vin := range admit.CoinsToRetain {
				output := topicInputs[topic][vin]
				if output == nil {
					return nil, ErrMissingInput
				}
				outputsConsumed = append(outputsConsumed, output)
				outpointsConsumed = append(outpointsConsumed, &output.Outpoint)
				delete(topicInputs[topic], vin)
			}
			for vin, output := range topicInputs[topic] {
				if err := e.deleteUTXODeep(ctx, output); err != nil {
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
				}
				if err := e.Storage.InsertOutput(ctx, output); err != nil {
					return nil, err
				}
				newOutpoints = append(newOutpoints, &output.Outpoint)
				for _, l := range e.LookupServices {
					if err := l.OutputAdded(ctx, output); err != nil {
						return nil, err
					}
				}
			}
			for _, output := range outputsConsumed {
				output.ConsumedBy = append(output.ConsumedBy, newOutpoints...)

				if err := e.Storage.UpdateConsumedBy(ctx, &output.Outpoint, output.Topic, output.ConsumedBy); err != nil {
					return nil, err
				}
			}
			if err := e.Storage.InsertAppliedTransaction(ctx, &overlay.AppliedTransaction{
				Txid:  txid,
				Topic: topic,
			}); err != nil {
				return nil, err
			}
		}
		if e.Advertiser == nil || mode == SubmitModeHistorical {
			return steak, nil
		}

		//TODO: Implement SYNC

		return steak, nil
	}
}

func (e *Engine) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	if l, ok := e.LookupServices[question.Service]; ok {
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
		if beef, err := tx.AtomicBEEF(false); err != nil {
			return nil, err
		} else {
			output.Beef = beef
			return output, nil
		}
	}
}

func (e *Engine) SyncAdvertisements() error {
	return nil
}

func (e *Engine) StartGASPSync() error {
	return nil
}

func (e *Engine) ProvideForeignSyncResponse(initialRequest *gasp.InitialRequest, topic string) (*gasp.InitialResponse, error) {
	return nil, nil
}

func (e *Engine) ProvideForeignGASPNode(graphId string, txid string, outputIndex uint32) (*gasp.GASPNode, error) {
	return nil, nil
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
			} else {
				e.updateInputProofs(ctx, input.SourceTransaction, txid, proof)
			}
		}
	}
	return nil
}

func (e *Engine) updateMerkleProof(ctx context.Context, output *Output, txid chainhash.Hash, proof *transaction.MerklePath) error {
	if len(output.Beef) == 0 {
		return errors.New("missing beef")
	} else if tx, err := transaction.NewTransactionFromBEEF(output.Beef); err != nil {
		return err
	} else if tx.MerklePath != nil {
		tx.MerklePath = proof
		return nil
	} else if err = e.updateInputProofs(ctx, tx, txid, proof); err != nil {
		return err
	} else {
		for _, outpoint := range output.ConsumedBy {
			if consumedOutputs, err := e.Storage.FindOutputsForTransaction(ctx, &outpoint.Txid, true); err != nil {
				return err
			} else {
				for _, consumed := range consumedOutputs {
					if err := e.updateMerkleProof(ctx, consumed, txid, proof); err != nil {
						return err
					}
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

func (e *Engine) GetDocumentationForLookupServiceProvider(provider string) (string, error) {
	if l, ok := e.LookupServices[provider]; !ok {
		return "", errors.New("no documentation found")
	} else {
		return l.GetDocumentation(), nil
	}
}

func (e *Engine) GetDocumentationForTopicManager(manager string) (string, error) {
	return "", nil
}

func FindPreviousTx(tx *transaction.Transaction, txid chainhash.Hash) *transaction.Transaction {
	if tx != nil {
		for _, input := range tx.Inputs {
			if input.SourceTXID.Equal(txid) {
				return input.SourceTransaction
			}
			if found := FindPreviousTx(input.SourceTransaction, *input.SourceTXID); found != nil {
				return found
			}
		}
	}
	return nil
}
