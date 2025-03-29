package engine

import (
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
)

type Output struct {
	Outpoint        overlay.Outpoint    `json:"-"`
	Topic           string              `json:"topic"`
	Script          *script.Script      `json:"-"`
	Satoshis        uint64              `json:"satoshis"`
	Spent           bool                `json:"spent"`
	OutputsConsumed []*overlay.Outpoint `json:""`
	ConsumedBy      []*overlay.Outpoint
	BlockHeight     uint32
	BlockIdx        uint64
	Beef            []byte
	AncillaryTxids  []*chainhash.Hash
	AncillaryBeef   []byte
}
