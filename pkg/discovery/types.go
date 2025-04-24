package discovery

import (
	"time"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

type SHIPRecord struct {
	Txid        *chainhash.Hash `json:"txid"`
	OutputIndex uint32          `json:"outputIndex"`
	IdentityKey *ec.PublicKey   `json:"identityKey"`
	Domain      string          `json:"domain"`
	Topic       string          `json:"topic"`
	CreatedAt   time.Time       `json:"createdAt"`
}

type SLAPRecord struct {
	Txid        *chainhash.Hash `json:"txid"`
	OutputIndex uint32          `json:"outputIndex"`
	IdentityKey *ec.PublicKey   `json:"identityKey"`
	Domain      string          `json:"domain"`
	Service     string          `json:"service"`
	CreatedAt   time.Time       `json:"createdAt"`
}

type SHIPQuery struct {
	Domain      string        `json:"domain,omitempty"`
	Topics      []string      `json:"topics,omitempty"`
	IdentityKey *ec.PublicKey `json:"identityKey,omitempty"`
}

type SLAPQuery struct {
	Domain      string        `json:"domain,omitempty"`
	Service     string        `json:"service,omitempty"`
	IdentityKey *ec.PublicKey `json:"identityKey,omitempty"`
}
