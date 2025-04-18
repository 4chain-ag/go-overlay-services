package engine

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/overlay"
)

type TopicManager interface {
	IdentifyAdmissableOutputs(ctx context.Context, beef []byte, previousCoins map[uint32][]byte) (overlay.AdmittanceInstructions, error)
	IdentifyNeededInputs(ctx context.Context, beef []byte) ([]*overlay.Outpoint, error)
	GetDocumentation() string
	GetMetaData() *overlay.MetaData
}
