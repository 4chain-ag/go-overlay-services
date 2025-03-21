package types

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

type TopicContext struct {
	Tx      *transaction.Transaction
	Inputs  map[uint32]*Output
	Outputs map[uint32]*Output
	Result  TopicResult
}

type TopicResult struct {
	Admit overlay.AdmittanceInstructions
}

type TopicManager interface {
	IdentifyAdmissableOutputs(tx *transaction.Transaction, loadInput func(uint32) (*Output, error)) (TopicResult, error)
	IdentifyNeededInputs(tx *transaction.Transaction) ([]*overlay.Outpoint, error)
	GetDependencies() []string
	GetDocumentation() string
	GetMetaData() *overlay.MetaData
}

type BaseTopicManager struct{}

func (b *BaseTopicManager) IdentifyAdmissableOutputs(tx *transaction.Transaction, loadInput func(uint32) error) (TopicResult, error) {
	return TopicResult{}, nil
}

func (b *BaseTopicManager) IdentifyNeededInputs(tx *transaction.Transaction) ([]*overlay.Outpoint, error) {
	return []*overlay.Outpoint{}, nil
}

func (b *BaseTopicManager) GetDocumentation() string {
	return ""
}
func (b *BaseTopicManager) GetMetaData() *overlay.MetaData {
	return &overlay.MetaData{}
}
