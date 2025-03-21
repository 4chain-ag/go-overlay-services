package engine

import (
	"bytes"
	"context"
	"errors"
	"slices"

	"github.com/4chain-ag/go-overlay-services/advertiser"
	"github.com/4chain-ag/go-overlay-services/gasp"
	"github.com/4chain-ag/go-overlay-services/storage"
	"github.com/4chain-ag/go-overlay-services/types"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
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

type Engine struct {
	Managers       map[string]types.TopicManager
	LookupServices map[string]types.LookupService
	Storage        storage.Storage
	ChainTracker   chaintracker.ChainTracker
	Broadcaster    transaction.Broadcaster
	SHIPTrackers   []string
	SLAPTrackers   []string
	Advertiser     *advertiser.Advertiser
}

var ErrUnknownTopic = errors.New("unknown-topic")
var ErrInvalidTransaction = errors.New("invalid-transaction")
var ErrMissingInput = errors.New("missing-input")

func (e *Engine) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode SumbitMode, onSteakReady func(steak overlay.Steak)) (map[string]*types.TopicContext, error) {
	for _, topic := range taggedBEEF.Topics {
		if _, ok := e.Managers[topic]; !ok {
			return nil, ErrUnknownTopic
		}
	}

	if tx, err := transaction.NewTransactionFromBEEF(taggedBEEF.Beef); err != nil {
		return nil, err
	} else if valid, err := spv.Verify(tx, e.ChainTracker, nil); err != nil {
		return nil, err
	} else if !valid {
		return nil, ErrInvalidTransaction
	} else {
		return e.submitTx(ctx, tx, taggedBEEF.Topics, mode, onSteakReady)
	}
}

func (e *Engine) submitTx(ctx context.Context, tx *transaction.Transaction, topics []string, mode SumbitMode, onSteakReady func(steak overlay.Steak)) (map[string]*types.TopicContext, error) {
	txid := tx.TxID()
	steak := make(overlay.Steak, len(topics))
	contexts := make(map[string]*types.TopicContext, len(topics))
	inpoints := make([]*overlay.Outpoint, 0, len(tx.Inputs))
	for _, input := range tx.Inputs {
		inpoints = append(inpoints, &overlay.Outpoint{
			Txid:        input.SourceTXID,
			OutputIndex: input.SourceTxOutIndex,
		})
	}
	for _, topic := range topics {
		manager := e.Managers[topic]
		if exists, err := e.Storage.DoesAppliedTransactionExist(ctx, &overlay.AppliedTransaction{
			Txid:  txid,
			Topic: topic,
		}); err != nil {
			return nil, err
		} else if exists {
			steak[topic] = &overlay.AdmittanceInstructions{}
			continue
		}
		tCtx := &types.TopicContext{
			Inputs:  make(map[uint32]*types.Output, len(tx.Inputs)),
			Outputs: make(map[uint32]*types.Output, len(tx.Outputs)),
		}
		contexts[topic] = tCtx
		var err error
		if tCtx.Result, err = manager.IdentifyAdmissableOutputs(tx, func(vin uint32) (output *types.Output, err error) {
			if vin >= uint32(len(tx.Inputs)) {
				return nil, ErrMissingInput
			}
			input := tx.Inputs[vin]
			outpoint := &overlay.Outpoint{
				Txid:        input.SourceTXID,
				OutputIndex: input.SourceTxOutIndex,
			}
			if output, err = e.Storage.FindOutput(ctx, outpoint, &topic, &FALSE, false); err != nil {
				return nil, err
			} else if output == nil {
				if input.SourceTransaction == nil {
					return nil, ErrMissingInput
				} else {
					var mode SumbitMode
					if input.SourceTransaction.MerklePath == nil {
						mode = SubmitModeCurrent
					} else {
						mode = SubmitModeHistorical
					}
					if tCtx, err := e.submitTx(ctx, input.SourceTransaction, []string{topic}, mode, nil); err != nil {
						return nil, err
					} else if output = tCtx[topic].Outputs[input.SourceTxOutIndex]; output == nil {
						out := input.SourceTransaction.Outputs[input.SourceTxOutIndex]
						output = &types.Output{
							Outpoint: outpoint,
							Script:   out.LockingScript,
							Satoshis: out.Satoshis,
						}
					}
				}
				tCtx.Inputs[vin] = output
			}
			return output, nil
		}); err != nil {
			return nil, err
		} else {
			steak[topic] = &tCtx.Result.Admit
		}
	}
	for _, topic := range topics {
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

	for _, topic := range topics {
		tCtx := contexts[topic]
		admittance := steak[topic]
		consumedOutpoints := make([]*overlay.Outpoint, 0, len(admittance.CoinsToRetain))
		consumedOutputs := make([]*types.Output, 0, len(admittance.CoinsToRetain))

		for vin, input := range tCtx.Inputs {
			if input == nil {
				continue
			}
			if !slices.Contains(admittance.CoinsToRetain, uint32(vin)) {
				admittance.CoinsRemoved = append(admittance.CoinsRemoved, uint32(vin))
				if err := e.deleteUTXODeep(ctx, input); err != nil {
					return nil, err
				}
			} else {
				if input := tCtx.Inputs[vin]; input == nil {
					return nil, ErrMissingInput
				} else {
					if tx.Inputs[vin].SourceTransaction == nil {
						return nil, ErrMissingInput
					}
					if err := e.Storage.InsertOutput(ctx, input); err != nil {
						return nil, err
					}
					for _, l := range e.LookupServices {
						if err := l.OutputAdded(ctx, input); err != nil {
							return nil, err
						}
					}
				}
				consumedOutpoints = append(consumedOutpoints, input.Outpoint)
				consumedOutputs = append(consumedOutputs, input)
			}
		}

		newOutpoints := make([]*overlay.Outpoint, 0, len(admittance.OutputsToAdmit))
		for _, vout := range admittance.OutputsToAdmit {
			out := tx.Outputs[vout]
			outpoint := &overlay.Outpoint{
				Txid:        txid,
				OutputIndex: uint32(vout),
			}
			output := &types.Output{
				Outpoint:        outpoint,
				Script:          out.LockingScript,
				Satoshis:        out.Satoshis,
				Topic:           topic,
				OutputsConsumed: consumedOutpoints,
			}
			tCtx.Outputs[uint32(vout)] = output
			if err := e.Storage.InsertOutput(ctx, output); err != nil {
				return nil, err
			}
			newOutpoints = append(newOutpoints, outpoint)
			for _, l := range e.LookupServices {
				if err := l.OutputAdded(ctx, output); err != nil {
					return nil, err
				}
			}
		}
		for _, output := range consumedOutputs {
			output.ConsumedBy = append(output.ConsumedBy, newOutpoints...)

			if err := e.Storage.UpdateConsumedBy(ctx, output.Outpoint, output.Topic, output.ConsumedBy); err != nil {
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
		return contexts, nil
	}

	//TODO: Implement SYNC

	return contexts, nil
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

func (e *Engine) GetUTXOHistory(ctx context.Context, output *types.Output, historySelector func(beef []byte, outputIndex uint32, currentDepth uint32) bool, currentDepth uint32) (*types.Output, error) {
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
	childHistories := make(map[string]*types.Output, len(outputsConsumed))
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
				Txid:        txin.SourceTXID,
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

func (e *Engine) deleteUTXODeep(ctx context.Context, output *types.Output) error {
	if len(output.ConsumedBy) == 0 {
		if err := e.Storage.DeleteOutput(ctx, output.Outpoint, output.Topic); err != nil {
			return err
		}
		for _, l := range e.LookupServices {
			if err := l.OutputDeleted(ctx, output.Outpoint, output.Topic); err != nil {
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
			if err := e.Storage.UpdateConsumedBy(ctx, staleOutput.Outpoint, staleOutput.Topic, staleOutput.ConsumedBy); err != nil {
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

func (e *Engine) updateMerkleProof(ctx context.Context, output *types.Output, txid chainhash.Hash, proof *transaction.MerklePath) error {
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
			if consumedOutputs, err := e.Storage.FindOutputsForTransaction(ctx, outpoint.Txid, true); err != nil {
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

func (e *Engine) GetDocumentationForLookupServiceProvider(provider string) string {
	if l, ok := e.LookupServices[provider]; ok {
		return "No documentation found!"
	} else {
		return l.GetDocumentation()
	}
}
