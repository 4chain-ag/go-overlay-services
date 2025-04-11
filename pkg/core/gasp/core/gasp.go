package core

import (
	"context"
	"fmt"
	"log"
	"slices"
	"sync"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

const MAX_CONCURRENCY = 16

type LogLevel int

var (
	LogLevelNone  LogLevel = 0
	LogLevelError LogLevel = 1
	LogLevelWarn  LogLevel = 2
	LogLevelInfo  LogLevel = 3
	LogLevelDebug LogLevel = 4
)

type GASPNodeRequest struct {
	GraphID     *overlay.Outpoint `json:"graphID"`
	Txid        *chainhash.Hash   `json:"txid"`
	OutputIndex uint32            `json:"outputIndex"`
	Metadata    bool              `json:"metadata"`
}

type GASPParams struct {
	Storage         GASPStorage
	Remote          GASPRemote
	LastInteraction uint32
	Version         *int
	LogPrefix       *string
	Unidirectional  bool
	LogLevel        *LogLevel
	Concurrency     int
}

type GASP struct {
	Version         int
	Remote          GASPRemote
	Storage         GASPStorage
	LastInteraction uint32
	LogPrefix       string
	Unidirectional  bool
	LogLevel        LogLevel
	limiter         chan struct{}
}

func NewGASP(params GASPParams) *GASP {
	gasp := &GASP{
		Storage:         params.Storage,
		Remote:          params.Remote,
		LastInteraction: params.LastInteraction,
		Unidirectional:  params.Unidirectional,
		// Sequential:      params.Sequential,
	}
	if params.Concurrency > 1 {
		gasp.limiter = make(chan struct{}, params.Concurrency)
	} else {
		gasp.limiter = make(chan struct{}, 1)
	}
	if params.Version != nil {
		gasp.Version = *params.Version
	} else {
		gasp.Version = 1
	}
	if params.LogPrefix != nil {
		gasp.LogPrefix = *params.LogPrefix
	} else {
		gasp.LogPrefix = "[GASP] "
	}
	if params.LogLevel != nil {
		gasp.LogLevel = *params.LogLevel
	} else {
		gasp.LogLevel = LogLevelInfo
	}
	return gasp
}

func (g *GASP) Sync(ctx context.Context) error {
	initialRequest := &GASPInitialRequest{
		Version: g.Version,
		Since:   g.LastInteraction,
	}
	initialResponse, err := g.Remote.GetInitialResponse(ctx, initialRequest)
	if err != nil {
		return err
	} else if len(initialResponse.UTXOList) > 0 {
		if foreignUTXOs, err := g.Storage.FindKnownUTXOs(ctx, 0); err != nil {
			return err
		} else {
			var wg sync.WaitGroup
			for _, outpoint := range initialResponse.UTXOList {
				if slices.ContainsFunc(foreignUTXOs, func(foreign *overlay.Outpoint) bool {
					return outpoint.Equal(foreign)
				}) {
					continue
				}
				wg.Add(1)
				g.limiter <- struct{}{}
				go func(outpoint *overlay.Outpoint) {
					defer func() {
						<-g.limiter
						wg.Done()
					}()
					var resolvedNode *GASPNode
					var err error
					if resolvedNode, err = g.Remote.RequestNode(ctx, outpoint, outpoint, true); err == nil {
						if err = g.processIncomingNode(ctx, resolvedNode, nil, &sync.Map{}); err == nil {
							if err = g.CompleteGraph(ctx, resolvedNode.GraphID); err == nil {
								return
							}
						}
					}
					log.Printf("Error with incoming UTXO %s: %v", outpoint.String(), err)
				}(outpoint)
			}
			wg.Wait()
		}
	}
	if !g.Unidirectional {
		if initialReply, err := g.Remote.GetInitialReplay(ctx, initialResponse); err != nil {
			return err
		} else {
			var wg sync.WaitGroup
			for _, outpoint := range initialReply.UTXOList {
				wg.Add(1)
				g.limiter <- struct{}{}
				go func(outpoint *overlay.Outpoint) {
					defer func() {
						<-g.limiter
						wg.Done()
					}()
					var outgoingNode *GASPNode
					if outgoingNode, err = g.Storage.HydrateGASPNode(ctx, outpoint, outpoint, true); err == nil {
						if err = g.processOutgoingNode(ctx, outgoingNode, &sync.Map{}); err == nil {
							return
						}
					}
					log.Printf("Error with outgoing UTXO %s: %v", outpoint, err)
				}(outpoint)
			}
			wg.Wait()
		}
	}
	return nil
}

func (g *GASP) GetInitialResponse2(ctx context.Context, request *GASPInitialRequest) (resp *GASPInitialResponse, err error) {
	if request.Version != g.Version {
		return nil, NewGASPVersionMismatchError(
			g.Version,
			request.Version,
		)
	}
	resp = &GASPInitialResponse{
		Since: g.LastInteraction,
	}
	if resp.UTXOList, err = g.Storage.FindKnownUTXOs(ctx, request.Since); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetInitialResponse is a method that handles the initial response from the GASP protocol.
// It checks the version of the request against the current version of the GASP instance.
// If the versions match, it retrieves the known UTXOs from the storage and returns them in the response.
// If the versions do not match, it returns a version mismatch error.
// The method also updates the LastInteraction field of the GASP instance with the current time.
// Categorize Feature: Wrap Error added for FindKnownUTXOs
func (g *GASP) GetInitialResponse(ctx context.Context, request *GASPInitialRequest) (resp *GASPInitialResponse, err error) {
	if request.Version != g.Version {
		return nil, NewGASPVersionMismatchError(
			g.Version,
			request.Version,
		)
	}

	resp = &GASPInitialResponse{
		Since: g.LastInteraction,
	}

	if resp.UTXOList, err = g.Storage.FindKnownUTXOs(ctx, request.Since); err != nil {
		return nil, fmt.Errorf("failed to find known UTXOs: %w", err)
	}

	return resp, nil
}

func (g *GASP) GetInitialReplay(ctx context.Context, response *GASPInitialResponse) (resp *GASPInitialReply, err error) {
	if knownUtxos, err := g.Storage.FindKnownUTXOs(ctx, response.Since); err != nil {
		return nil, err
	} else {
		resp = &GASPInitialReply{
			UTXOList: make([]*overlay.Outpoint, 0, len(response.UTXOList)),
		}
		for _, utxo := range response.UTXOList {
			if !slices.ContainsFunc(knownUtxos, func(known *overlay.Outpoint) bool {
				return known.Equal(utxo)
			}) {
				resp.UTXOList = append(resp.UTXOList, utxo)
			}
		}
		return resp, nil
	}
}

func (g *GASP) RequestNode(ctx context.Context, graphID *overlay.Outpoint, outpoint *overlay.Outpoint, metadata bool) (node *GASPNode, err error) {
	if node, err = g.Storage.HydrateGASPNode(ctx, graphID, outpoint, metadata); err != nil {
		return nil, err
	}
	return node, nil
}

func (g *GASP) SubmitNode(ctx context.Context, node *GASPNode) (requestedInputs *GASPNodeResponse, err error) {
	if err = g.Storage.AppendToGraph(ctx, node, nil); err != nil {
		return nil, err
	} else if requestedInputs, err = g.Storage.FindNeededInputs(ctx, node); err != nil {
		return nil, err
	} else if requestedInputs != nil {
		if err := g.CompleteGraph(ctx, node.GraphID); err != nil {
			return nil, err
		}
	}
	return requestedInputs, nil
}

func (g *GASP) CompleteGraph(ctx context.Context, graphID *overlay.Outpoint) (err error) {
	if err = g.Storage.ValidateGraphAnchor(ctx, graphID); err == nil {
		if err := g.Storage.FinalizeGraph(ctx, graphID); err == nil {
			return nil
		}
	}
	log.Printf("Error completing graph %s: %v", graphID.String(), err)
	return g.Storage.DiscardGraph(ctx, graphID)
}

func (g *GASP) processIncomingNode(ctx context.Context, node *GASPNode, spentBy *chainhash.Hash, seenNodes *sync.Map) error {
	if txid, err := g.computeTxID(node.RawTx); err != nil {
		return err
	} else {
		nodeId := (&overlay.Outpoint{
			Txid:        *txid,
			OutputIndex: node.OutputIndex,
		}).String()
		if _, ok := seenNodes.Load(nodeId); ok {
			return nil
		}
		seenNodes.Store(nodeId, struct{}{})
		if err := g.Storage.AppendToGraph(ctx, node, spentBy); err != nil {
			return err
		} else if neededInputs, err := g.Storage.FindNeededInputs(ctx, node); err != nil {
			return err
		} else if neededInputs != nil {
			var wg sync.WaitGroup
			errors := make(chan error)
			for outpointStr, data := range neededInputs.RequestedInputs {
				wg.Add(1)
				g.limiter <- struct{}{}
				go func(outpointStr string, data *GASPNodeResponseData) {
					defer func() {
						<-g.limiter
						wg.Done()
					}()
					if outpoint, err := overlay.NewOutpointFromString(outpointStr); err != nil {
						errors <- err
					} else if newNode, err := g.Remote.RequestNode(ctx, node.GraphID, outpoint, data.Metadata); err != nil {
						errors <- err
					} else if err := g.processIncomingNode(ctx, newNode, txid, seenNodes); err != nil {
						errors <- err
					}
				}(outpointStr, data)
			}
			go func() {
				wg.Wait()
				close(errors)
			}()
			for err := range errors {
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (g *GASP) processOutgoingNode(ctx context.Context, node *GASPNode, seenNodes *sync.Map) error {
	if g.Unidirectional {
		return nil
	}
	if txid, err := g.computeTxID(node.RawTx); err != nil {
		return err
	} else {
		nodeId := (&overlay.Outpoint{
			Txid:        *txid,
			OutputIndex: node.OutputIndex,
		}).String()
		if _, ok := seenNodes.Load(nodeId); ok {
			return nil
		}
		seenNodes.Store(nodeId, struct{}{})
		if response, err := g.Remote.SubmitNode(ctx, node); err != nil {
			return err
		} else if response != nil {
			var wg sync.WaitGroup
			for outpointStr, data := range response.RequestedInputs {
				wg.Add(1)
				g.limiter <- struct{}{}
				go func(outpointStr string, data *GASPNodeResponseData) {
					defer func() {
						<-g.limiter
						wg.Done()
					}()
					var outpoint *overlay.Outpoint
					var err error
					if outpoint, err = overlay.NewOutpointFromString(outpointStr); err == nil {
						var hydratedNode *GASPNode
						if hydratedNode, err = g.Remote.RequestNode(ctx, node.GraphID, outpoint, data.Metadata); err == nil {
							if err = g.processOutgoingNode(ctx, hydratedNode, seenNodes); err == nil {
								return
							}
						}
					}
					log.Printf("Error hydrating node: %v", err)
				}(outpointStr, data)
			}
			wg.Wait()
		}
	}
	return nil
}

func (g *GASP) computeTxID(rawtx string) (*chainhash.Hash, error) {
	if tx, err := transaction.NewTransactionFromHex(rawtx); err != nil {
		return nil, err
	} else {
		return tx.TxID(), nil
	}
}
