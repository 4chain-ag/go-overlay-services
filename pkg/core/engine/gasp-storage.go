package engine

import (
	"bytes"
	"context"
	"errors"
	"slices"
	"sync"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/spv"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

var ErrGraphFull = errors.New("graph is full")

type GraphNode struct {
	core.GASPNode
	Txid     *chainhash.Hash `json:"txid"`
	SpentBy  *chainhash.Hash `json:"spentBy"`
	Children []*GraphNode    `json:"children"`
	Parent   *GraphNode      `json:"parent"`
}

type OverlayGASPStorage struct {
	Topic              string
	Engine             *Engine
	MaxNodesInGraph    *int
	tempGraphNodeRefs  sync.Map
	tempGraphNodeCount int
}

func NewOverlayGASPStorage(topic string, engine *Engine, maxNodesInGraph *int) *OverlayGASPStorage {
	return &OverlayGASPStorage{
		Topic:           topic,
		Engine:          engine,
		MaxNodesInGraph: maxNodesInGraph,
	}
}

func (s *OverlayGASPStorage) FindKnownUTXOs(ctx context.Context, since uint32) ([]*overlay.Outpoint, error) {
	if utxos, err := s.Engine.Storage.FindUTXOsForTopic(ctx, s.Topic, since, false); err != nil {
		return nil, err
	} else {
		outpoints := make([]*overlay.Outpoint, len(utxos))
		for i, utxo := range utxos {
			outpoints[i] = &utxo.Outpoint
		}

		return outpoints, nil
	}
}

func (s *OverlayGASPStorage) HydrateGASPNode(ctx context.Context, graphID *overlay.Outpoint, outpoint *overlay.Outpoint, metadata bool) (*core.GASPNode, error) {
	if output, err := s.Engine.Storage.FindOutput(ctx, outpoint, nil, nil, true); err != nil {
		return nil, err
	} else if output.Beef == nil {
		return nil, ErrMissingInput
	} else if tx, err := transaction.NewTransactionFromBEEF(output.Beef); err != nil {
		return nil, err
	} else {
		node := &core.GASPNode{
			GraphID:     graphID,
			OutputIndex: outpoint.OutputIndex,
			RawTx:       tx.Hex(),
		}
		if tx.MerklePath != nil {
			proof := tx.MerklePath.Hex()
			node.Proof = &proof
		}
		return node, nil
	}
}

func (s *OverlayGASPStorage) FindNeededInputs(ctx context.Context, gaspTx *core.GASPNode) (*core.GASPNodeResponse, error) {
	response := &core.GASPNodeResponse{
		RequestedInputs: make(map[string]*core.GASPNodeResponseData),
	}
	tx, err := transaction.NewTransactionFromHex(gaspTx.RawTx)
	if err != nil {
		return nil, err
	}
	if gaspTx.Proof == nil {
		for _, input := range tx.Inputs {
			outpoint := &overlay.Outpoint{
				Txid:        *input.SourceTXID,
				OutputIndex: input.SourceTxOutIndex,
			}
			response.RequestedInputs[outpoint.String()] = &core.GASPNodeResponseData{
				Metadata: false,
			}
		}

		return s.stripAlreadyKnowInputs(ctx, response)
	} else if tx.MerklePath, err = transaction.NewMerklePathFromHex(*gaspTx.Proof); err != nil {
		return nil, err
	}
	if beef, err := transaction.NewBeefFromTransaction(tx); err != nil {
		return nil, err
	} else {
		if len(gaspTx.AncillaryBeef) > 0 {
			if err := beef.MergeBeefBytes(gaspTx.AncillaryBeef); err != nil {
				return nil, err
			}
		}
		previousCoins := make(map[uint32][]byte, len(tx.Inputs))
		for vin, input := range tx.Inputs {
			outpoint := &overlay.Outpoint{
				Txid:        *input.SourceTXID,
				OutputIndex: input.SourceTxOutIndex,
			}
			if output, err := s.Engine.Storage.FindOutput(ctx, outpoint, &s.Topic, nil, true); err != nil {
				return nil, err
			} else if output != nil {
				previousCoins[uint32(vin)] = output.Beef
			}
		}

		if beefBytes, err := beef.AtomicBytes(tx.TxID()); err != nil {
			return nil, err
		} else if admit, err := s.Engine.Managers[s.Topic].IdentifyAdmissableOutputs(ctx, beefBytes, previousCoins); err != nil {
			return nil, err
		} else if !slices.Contains(admit.OutputsToAdmit, gaspTx.OutputIndex) {
			if neededInputs, err := s.Engine.Managers[s.Topic].IdentifyNeededInputs(ctx, beefBytes); err != nil {
				return nil, err
			} else {
				for _, outpoint := range neededInputs {
					response.RequestedInputs[outpoint.String()] = &core.GASPNodeResponseData{
						Metadata: true,
					}
				}
				return s.stripAlreadyKnowInputs(ctx, response)
			}
		}
	}

	return response, nil
}

func (s *OverlayGASPStorage) stripAlreadyKnowInputs(ctx context.Context, response *core.GASPNodeResponse) (*core.GASPNodeResponse, error) {
	if response == nil {
		return nil, nil
	}
	for outpointStr := range response.RequestedInputs {
		if outpoint, err := overlay.NewOutpointFromString(outpointStr); err != nil {
			return nil, err
		} else if found, err := s.Engine.Storage.FindOutput(ctx, outpoint, &s.Topic, nil, false); err != nil {
			return nil, err
		} else if found != nil {
			delete(response.RequestedInputs, outpointStr)
		}
	}
	if len(response.RequestedInputs) == 0 {
		return nil, nil
	}
	return response, nil
}

func (s *OverlayGASPStorage) AppendToGraph(ctx context.Context, gaspTx *core.GASPNode, spentBy *chainhash.Hash) error {
	if s.MaxNodesInGraph != nil && s.tempGraphNodeCount >= *s.MaxNodesInGraph {
		return ErrGraphFull
	}

	if tx, err := transaction.NewTransactionFromHex(gaspTx.RawTx); err != nil {
		return err
	} else {
		txid := tx.TxID()
		if gaspTx.Proof != nil {
			if tx.MerklePath, err = transaction.NewMerklePathFromHex(*gaspTx.Proof); err != nil {
				return err
			}
		}
		newGraphNode := &GraphNode{
			GASPNode: *gaspTx,
			Txid:     txid,
			Children: []*GraphNode{},
		}
		if spentBy == nil {
			if _, ok := s.tempGraphNodeRefs.LoadOrStore(gaspTx.GraphID.String(), newGraphNode); !ok {
				s.tempGraphNodeCount++
			}
		} else if node, ok := s.tempGraphNodeRefs.Load(gaspTx.GraphID.String()); ok {
			parentNode := node.(*GraphNode)
			parentNode.Children = append(parentNode.Children, newGraphNode)
			newGraphNode.Parent = parentNode
			newGraphOutpoint := &overlay.Outpoint{
				Txid:        *txid,
				OutputIndex: gaspTx.OutputIndex,
			}
			if _, ok := s.tempGraphNodeRefs.LoadOrStore(newGraphOutpoint.String(), newGraphNode); !ok {
				s.tempGraphNodeCount++
			}
		} else {
			return ErrMissingInput
		}
		return nil
	}
}

func (s *OverlayGASPStorage) ValidateGraphAnchor(ctx context.Context, graphID *overlay.Outpoint) error {
	if rootNode, ok := s.tempGraphNodeRefs.Load(graphID.String()); !ok {
		return ErrMissingInput
	} else if beef, err := s.getBEEFForNode(rootNode.(*GraphNode)); err != nil {
		return err
	} else if tx, err := transaction.NewTransactionFromBEEF(beef); err != nil {
		return err
	} else if valid, err := spv.Verify(tx, s.Engine.ChainTracker, nil); err != nil {
		return err
	} else if !valid {
		return errors.New("graph anchor is not a valid transaction")
	} else if beefs, err := s.computeOrderedBEEFsForGraph(ctx, graphID); err != nil {
		return err
	} else {
		coins := make(map[string]struct{})
		for _, beefBytes := range beefs {
			if tx, err := transaction.NewTransactionFromBEEF(beefBytes); err != nil {
				return err
			} else {
				previousCoins := make(map[uint32][]byte, len(tx.Inputs))
				for vin, input := range tx.Inputs {
					outpoint := &overlay.Outpoint{
						Txid:        *input.SourceTXID,
						OutputIndex: input.SourceTxOutIndex,
					}
					if output, err := s.Engine.Storage.FindOutput(ctx, outpoint, &s.Topic, nil, true); err != nil {
						return err
					} else if output != nil {
						previousCoins[uint32(vin)] = output.Beef
					}
				}
				if admit, err := s.Engine.Managers[s.Topic].IdentifyAdmissableOutputs(ctx, beef, previousCoins); err != nil {
					return err
				} else {
					for _, vout := range admit.OutputsToAdmit {
						outpoint := &overlay.Outpoint{
							Txid:        *tx.TxID(),
							OutputIndex: vout,
						}
						coins[outpoint.String()] = struct{}{}
					}
				}
			}
		}
		if _, ok := coins[graphID.String()]; !ok {
			return errors.New("graph did not result in topical admittance of the root node. rejecting")
		}
		return nil
	}
}

func (s *OverlayGASPStorage) DiscardGraph(ctx context.Context, graphID *overlay.Outpoint) error {
	s.tempGraphNodeRefs.Range(func(nodeId, graphRef any) bool {
		if graphRef.(*GraphNode).GraphID.Equal(graphID) {
			s.tempGraphNodeRefs.Delete(nodeId)
			s.tempGraphNodeCount--
			return false
		}
		return true
	})
	return nil
}

func (s *OverlayGASPStorage) FinalizeGraph(ctx context.Context, graphID *overlay.Outpoint) error {
	if beefs, err := s.computeOrderedBEEFsForGraph(ctx, graphID); err != nil {
		return err
	} else {
		for _, beef := range beefs {
			if _, err := s.Engine.Submit(
				ctx,
				overlay.TaggedBEEF{
					Topics: []string{s.Topic},
					Beef:   beef,
				},
				SubmitModeHistorical,
				nil,
			); err != nil {
				return err
			}
		}
		return nil
	}
}

func (s *OverlayGASPStorage) computeOrderedBEEFsForGraph(ctx context.Context, graphID *overlay.Outpoint) ([][]byte, error) {
	beefs := make([][]byte, 0)
	var hydrator func(node *GraphNode) error
	hydrator = func(node *GraphNode) error {
		if currentBeef, err := s.getBEEFForNode(node); err != nil {
			return err
		} else {
			if slices.IndexFunc(beefs, func(beef []byte) bool {
				return bytes.Equal(beef, currentBeef)
			}) == -1 {
				beefs = append([][]byte{currentBeef}, beefs...)
			}
			for _, child := range node.Children {
				if err := hydrator(child); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if foundRoot, ok := s.tempGraphNodeRefs.Load(graphID.String()); !ok {
		return nil, errors.New("unable to find root node in graph for finalization")
	} else if err := hydrator(foundRoot.(*GraphNode)); err != nil {
		return nil, err
	} else {
		return beefs, nil
	}
}

func (s *OverlayGASPStorage) getBEEFForNode(node *GraphNode) ([]byte, error) {
	var hydrator func(node *GraphNode) (*transaction.Transaction, error)
	hydrator = func(node *GraphNode) (*transaction.Transaction, error) {
		if tx, err := transaction.NewTransactionFromHex(node.RawTx); err != nil {
			return nil, err
		} else if node.Proof != nil {
			if tx.MerklePath, err = transaction.NewMerklePathFromHex(*node.Proof); err != nil {
				return nil, err
			}
			return tx, nil
		} else {
			for vin, input := range tx.Inputs {
				outpoint := &overlay.Outpoint{
					Txid:        *input.SourceTXID,
					OutputIndex: input.SourceTxOutIndex,
				}
				if foundNode, ok := s.tempGraphNodeRefs.Load(outpoint.String()); !ok {
					return nil, errors.New("required input node for unproven parent not found in temporary graph store")
				} else if tx.Inputs[vin].SourceTransaction, err = hydrator(foundNode.(*GraphNode)); err != nil {
					return nil, err
				}
			}
			return tx, nil
		}
	}
	if tx, err := hydrator(node); err != nil {
		return nil, err
	} else if beef, err := transaction.NewBeefFromTransaction(tx); err != nil {
		return nil, err
	} else {
		if len(node.AncillaryBeef) > 0 {
			if err := beef.MergeBeefBytes(node.AncillaryBeef); err != nil {
				return nil, err
			}
		}
		return beef.AtomicBytes(tx.TxID())
	}
}
