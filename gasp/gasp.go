package gasp

import "github.com/bsv-blockchain/go-sdk/overlay"

type InitialRequest struct {
	Version int    `json:"version"`
	Since   uint32 `json:"since"`
}

type InitialResponse struct {
	UTXOList []*overlay.Outpoint `json:"UTXOList"`
	Since    uint32              `json:"since"`
}

type InitialReply struct {
	UTXOList []*overlay.Outpoint `json:"UTXOList"`
}

type GASPInput struct {
	Hash string `json:"hash"`
}

type GASPNode struct {
	GraphID        string                `json:"graphID"`
	RawTx          string                `json:"rawTx"`
	OutputIndex    uint32                `json:"outputIndex"`
	Proof          string                `json:"proof"`
	TxMetadata     string                `json:"txMetadata"`
	OutputMetadata string                `json:"outputMetadata"`
	Inputs         map[string]*GASPInput `json:"inputs"`
}
