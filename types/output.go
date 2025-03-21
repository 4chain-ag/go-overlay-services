package types

import (
	"encoding/json"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
)

type Output struct {
	Outpoint        *overlay.Outpoint   `json:"-"`
	Topic           string              `json:"topic"`
	Script          *script.Script      `json:"-"`
	Satoshis        uint64              `json:"satoshis"`
	Spent           bool                `json:"spent"`
	OutputsConsumed []*overlay.Outpoint `json:""`
	ConsumedBy      []*overlay.Outpoint
	BlockHeight     uint32
	BlockIdx        uint64
	Beef            []byte
}

func (o *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"txid":             o.Outpoint.Txid.String(),
		"vout":             o.Outpoint.OutputIndex,
		"height":           o.BlockHeight,
		"idx":              o.BlockIdx,
		"satoshis":         o.Satoshis,
		"script":           o.Script.String(),
		"spent":            o.Spent,
		"outputs_consumed": o.OutputsConsumed,
		"consumed_by":      o.ConsumedBy,
		"topic":            o.Topic,
		"beef":             o.Beef,
	})
}
