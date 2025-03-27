package dto

import (
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
)

type Output struct {
	Outpoint        overlay.Outpoint
	Topic           string
	Script          *script.Script
	Satoshis        uint64
	Spent           bool
	OutputsConsumed []*overlay.Outpoint
	ConsumedBy      []*overlay.Outpoint
	BlockHeight     uint32
	BlockIdx        uint64
	Beef            []byte
	Dependencies    []*chainhash.Hash
}

type InsertOutput struct {
	TxID        string
	Vout        uint32
	Topic       string
	BlockHeight uint32
	BlockIdx    uint64
	Satoshis    uint64
	Script      []byte
	Spent       bool
}

type FindOutput struct {
	TxID        string
	OutputIndex uint32
	Topic       string
	Spent       bool
	IncludeBEEF bool
}
