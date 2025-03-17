package engine

import (
	"bytes"
	"context"
	"fmt"
	"slices"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/bsv-blockchain/go-sdk/overlay/topic"
	"github.com/bsv-blockchain/go-sdk/spv"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/chaintracker"

	"github.com/4chain-ag/go-overlay-services/advertiser"
	"github.com/4chain-ag/go-overlay-services/storage"
)

type SumbitMode string

var (
	SubmitModeHistorical SumbitMode = "historical-tx"
	SubmitModeCurrent    SumbitMode = "current-tx"
)

type Engine struct {
	Managers       map[string]topic.TopicManager
	LookupServices map[string]lookup.LookupService
	Storage        storage.Storage
	ChainTracker   chaintracker.ChainTracker
	Broadcaster    transaction.Broadcaster
	SHIPTrackers   []string
	SLAPTrackers   []string
	Advertiser     *advertiser.Advertiser
}

var ErrUnknownTopic = fmt.Errorf("unknown-topic")
var ErrInvalidTransaction = fmt.Errorf("invalid-transaction")

func (e *Engine) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode SumbitMode, onCtxReady func(subCtx *overlay.SubmitContext)) (*overlay.SubmitContext, error) {
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
		subCtx := &overlay.SubmitContext{
			Txid:            tx.TxID(),
			Tx:              tx,
			Beef:            taggedBEEF.Beef,
			TopicAdmittance: make(overlay.Steak),
			TopicInputs:     make(map[string][]*overlay.Output),
		}
		for _, topic := range taggedBEEF.Topics {
			if err := e.ExecuteTopicManager(ctx, subCtx, topic); err != nil {
				return nil, err
			}
		}
		if mode != SubmitModeHistorical && e.Broadcaster != nil {
			if _, failure := e.Broadcaster.Broadcast(subCtx.Tx); failure != nil {
				return nil, failure
			}
		}

		if onCtxReady != nil {
			onCtxReady(subCtx)
		}

		for _, topic := range taggedBEEF.Topics {
			if exists, err := e.Storage.DoesAppliedTransactionExist(ctx, &overlay.AppliedTransaction{
				Txid:  subCtx.Txid,
				Topic: topic,
			}); err != nil {
				return nil, err
			} else if exists {
				subCtx.TopicAdmittance[topic] = &overlay.Admittance{}
				continue
			}
			admittance := subCtx.TopicAdmittance[topic]
			inputs := subCtx.TopicInputs[topic]
			consumedOutpoints := make([]*overlay.Outpoint, 0, len(admittance.CoinsToRetain))
			consumedOutputs := make([]*overlay.Output, 0, len(admittance.CoinsToRetain))

			for vin, input := range inputs {
				if input == nil {
					continue
				}
				if slices.Contains(admittance.CoinsToRetain, uint32(vin)) {
					consumedOutpoints = append(consumedOutpoints, input.Outpoint)
					consumedOutputs = append(consumedOutputs, input)
				} else {
					admittance.CoinsRemoved = append(admittance.CoinsRemoved, uint32(vin))
					if err := e.deleteUTXODeep(ctx, input); err != nil {
						return nil, err
					}
				}
			}

			newOutpoints := make([]*overlay.Outpoint, 0, len(admittance.OutputsToAdmit))
			for _, vout := range admittance.OutputsToAdmit {
				out := subCtx.Tx.Outputs[vout]
				outpoint := &overlay.Outpoint{
					Txid: subCtx.Txid,
					Vout: uint32(vout),
				}
				if err := e.Storage.InsertOutput(ctx, &overlay.Output{
					Outpoint:        outpoint,
					Script:          *out.LockingScript,
					Satoshis:        out.Satoshis,
					Topic:           topic,
					OutputsConsumed: consumedOutpoints,
				}); err != nil {
					return nil, err
				}
				newOutpoints = append(newOutpoints, outpoint)
				for _, l := range e.LookupServices {
					if err := l.OutputAdded(*subCtx, vout, topic); err != nil {
						return nil, err
					}
				}
			}
			for _, output := range consumedOutputs {
				outpointSet := make(map[string]struct{})
				consumedBy := make([]*overlay.Outpoint, 0, len(output.ConsumedBy)+len(newOutpoints))
				for _, outpoint := range output.ConsumedBy {
					op := outpoint.String()
					if _, ok := outpointSet[op]; !ok {
						consumedBy = append(consumedBy, outpoint)
						outpointSet[op] = struct{}{}
					}
				}
				for _, outpoint := range newOutpoints {
					op := outpoint.String()
					if _, ok := outpointSet[op]; !ok {
						consumedBy = append(consumedBy, outpoint)
						outpointSet[op] = struct{}{}
					}
				}
				if err := e.Storage.UpdateConsumedBy(ctx, output.Outpoint, output.Topic, consumedBy); err != nil {
					return nil, err
				}
			}
			if err := e.Storage.InsertAppliedTransaction(ctx, &overlay.AppliedTransaction{
				Txid:  subCtx.Txid,
				Topic: topic,
			}); err != nil {
				return nil, err
			}
		}
		if e.Advertiser == nil || mode == SubmitModeHistorical {
			return subCtx, nil
		}

		//TODO: Implement SYNC

		return subCtx, nil
	}
}

func (e *Engine) ExecuteTopicManager(ctx context.Context, subCtx *overlay.SubmitContext, topic string) (err error) {
	if _, ok := subCtx.TopicAdmittance[topic]; ok {
		return nil
	}
	manager := e.Managers[topic]
	deps := manager.GetDependencies()
	for _, dep := range deps {
		if err := e.ExecuteTopicManager(ctx, subCtx, dep); err != nil {
			return err
		}
	}
	outpoints := make([]*overlay.Outpoint, 0, len(subCtx.Tx.Inputs))
	for _, input := range subCtx.Tx.Inputs {
		outpoints = append(outpoints, &overlay.Outpoint{
			Txid: input.SourceTXID,
			Vout: input.SourceTxOutIndex,
		})
	}
	if subCtx.TopicInputs[topic], err = e.Storage.FindOutputs(ctx, outpoints, topic, false, false); err != nil {
		return err
	} else if subCtx.TopicAdmittance[topic], err = manager.IdentifyAdmissableOutputs(*subCtx); err != nil {
		return err
	} else if err := e.Storage.MarkUTXOsAsSpent(ctx, outpoints, topic, subCtx.Txid); err != nil {
		return err
	}

	return nil
}

func (e *Engine) deleteUTXODeep(ctx context.Context, output *overlay.Output) error {
	if len(output.ConsumedBy) == 0 {
		if err := e.Storage.DeleteOutput(ctx, output.Outpoint, output.Topic); err != nil {
			return err
		}
		for _, l := range e.LookupServices {
			if err := l.OutputDeleted(output.Outpoint, output.Topic); err != nil {
				return err
			}
		}
	}
	if len(output.OutputsConsumed) == 0 {
		return nil
	}

	for _, outpoint := range output.OutputsConsumed {
		staleOutput, err := e.Storage.FindOutput(ctx, outpoint, output.Topic, false, false)
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
